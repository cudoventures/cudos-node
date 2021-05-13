const os = require('os');
const fs = require('fs');
const asyncFs = require('fs/promises');
const path = require('path');
const { exec } = require('child_process');

const fse = require('fs-extra');
const { ArgumentParser } = require('argparse');
const archiver = require('archiver');
const SCP2 = require('scp2');
const SSH2Client = require('ssh2').Client;

const SecretsConfig = require('./secrets.json');

const TEMP_DIR = path.join(os.tmpdir(), 'cudos-builder');

const LOCAL_BUILD = true;

async function main() {
    const args = getArgParser();
    const secrets = getSecrets(args.target);

    console.log(`Temp dir: ${TEMP_DIR}`);
    console.log(`Target: ${args.target}`);
    console.log(`Server Dir: ${secrets.serverPath}`);

    await asyncFs.rm(TEMP_DIR, { 'recursive': true, 'force': true });
    const { deployFilePath, deployFilename } = await initTempDirectory();

    try {
        if (LOCAL_BUILD === true) {
            console.log('Building meteor');
            await buildMeteor();
        }

        console.log('Creating archive');
        await createArchive(deployFilePath, deployFilename);

        console.log('Uploading');
        await uploadFile(secrets, deployFilePath, deployFilename);

        console.log('Execute commands');
        await executeCommands(args, secrets, deployFilePath, deployFilename);
    } finally {
        try {
            // await asyncFs.access(deployFilePath, fs.constants.F_OK);
            await asyncFs.rm(TEMP_DIR, { 'recursive': true, 'force': true });
        } catch (e) {
        }
    }
}

function getArgParser() {
    const parser = new ArgumentParser({description: 'Cudos testnet root node deployer'});
    parser.add_argument('--target', { 'required': true, 'choices': ['testnet'] });
    // parser.add_argument('--init', { 'required': true, 'choices': ['0', '1'] });
    return parser.parse_args();
}

function getSecrets(target) {
    const secrets = SecretsConfig[target];
    if (secrets === undefined) {
        console.error(`Secrets with target not found. Target:${target}`);
        return;
    }

    return secrets
}

async function initTempDirectory() {
    return new Promise(async (resolve, reject) => {
        try {
            await asyncFs.access(TEMP_DIR);
        } catch (e) {
            await asyncFs.mkdir(TEMP_DIR, { 'recursive': true });
        }

        const deployFilename = `cudos-utils.tar.gz`;
        const deployFilePath = path.join(TEMP_DIR, deployFilename);
        resolve({
            deployFilePath,
            deployFilename,
        })
    });
}

async function buildMeteor() {
    const source = path.resolve('../project-explorer');
    const target = path.join(TEMP_DIR, 'project-explorer');
    const copyFilter = [
        path.resolve(path.join(source, 'node_modules')),
        path.resolve(path.join(source, '.meteor/local')),
    ];

    await fse.copy(source, target, {
        filter: (src, desc) => {
            for (let i = copyFilter.length; i-- > 0; ) {
                if (src.startsWith(copyFilter[i])) {
                    return false;
                }
            }
            return true;
        }
    });

    await new Promise((resolve, reject) => {
        const cmd = [
            `cd ${target}`,
            'npm i',
            'meteor build ../meteor-build/ --architecture os.linux.x86_64 --server-only --allow-superuser',
            'cd ../meteor-build',
            'tar -zxf ./project-explorer.tar.gz',
            'cp ../project-explorer/run-testnet.sh ./bundle/run-testnet.sh',
        ];

        exec(cmd.join('&&'), (err, stdout, etderr) => {
            if (err !== null) {
                reject(err);
            }

            resolve();
        });
    });
}

async function createArchive(deployFilePath, deployFilename) {
    return new Promise(async (resolve, reject) => {
        const output = fs.createWriteStream(deployFilePath);
        const archive = archiver('zip', {
            zlib: { level: 9 }, // Sets the compression level.
        });

        output.on('close', () => {
            console.log(`${archive.pointer()} total bytes`);
            console.log('archiver has been finalized and the output file descriptor has closed.');
            resolve({
                deployFilePath,
                deployFilename,
            });
        });

        output.on('end', () => {
            console.log('Data has been drained');
        });

        archive.on('warning', (err) => {
            if (err.code === 'ENOENT') {
                // log warning
            } else {
                // throw error
                reject(err);
            }
        });

        // good practice to catch this error explicitly
        archive.on('error', (err) => {
            reject(err);
        });

        // pipe archive data to the file
        archive.pipe(output);

        // append files from a sub-directory, putting its contents at the root of archive
        archive.directory(path.resolve('../project-faucet-cli'), "/project-faucet-cli");
        if (LOCAL_BUILD === true) {
            archive.directory(path.resolve(path.join(TEMP_DIR, 'meteor-build/bundle')), "/project-explorer");
        }

        if (LOCAL_BUILD === false) {
            const projectExplorerAbsPath = path.resolve('../project-explorer');
            const pathContent = await asyncFs.readdir(projectExplorerAbsPath);
            for (let i = 0;  i < pathContent.length; ++i) {
                const itemAbsPath = path.join(projectExplorerAbsPath, pathContent[i]);
                const stat = await asyncFs.stat(itemAbsPath);

                switch (pathContent[i]) {
                    case 'node_modules':
                        break;
                    case '.meteor':
                        const meteorPathContent = await asyncFs.readdir(itemAbsPath);
                        for (let j = 0;  j < meteorPathContent.length;  ++j) {
                            if (meteorPathContent[j] === 'local') {
                                continue;
                            }
                            const meteorItemAbsPath = path.join(projectExplorerAbsPath, pathContent[i], meteorPathContent[j]);
                            const meteorStat = await asyncFs.stat(meteorItemAbsPath);
                            if (meteorStat.isDirectory() === true) {
                                archive.directory(meteorItemAbsPath, `/project-explorer/${pathContent[i]}/${meteorPathContent[j]}`);
                            } else {
                                archive.file(meteorItemAbsPath, { 'name': `/project-explorer/${pathContent[i]}/${meteorPathContent[j]}` } );
                            }
                        }
                        break;
                    default:
                        if (stat.isDirectory() === true) {
                            archive.directory(itemAbsPath, `/project-explorer/${pathContent[i]}`);
                        } else {
                            archive.file(itemAbsPath, { 'name': `/project-explorer/${pathContent[i]}` } );
                        }
                        break;
                }
            }
        }

        archive.finalize();
    });
}

function uploadFile(secrets, deployFilePath, deployFilename) {
    return new Promise(async (resolve, reject) => {
        const spcClient = new SCP2.Client();
        spcClient.on('connect', () => {
            console.log('Connected to server.');
        });

        spcClient.on('transfer', (buffer, uploaded, total) => {
            console.log(`Uploaded: ${uploaded + 1}/${total}`);
        });

        const destOptions = {
            host: secrets.host,
            port: secrets.port,
            username: secrets.username,
            passphrase: secrets.keyPass,
            privateKey: (await asyncFs.readFile(secrets.privateKey)).toString(),
            path: secrets.serverPath,
        };

        SCP2.scp(deployFilePath, destOptions, spcClient, (err) => {
            if (err) {
                console.error('Error:', err);
                reject(err);
                return;
            }

            resolve();
        });
    });
}

async function executeCommands(args, secrets, deployFilePath, deployFilename) {
    const conn = new SSH2Client();
    const filePath = path.join(secrets.serverPath, deployFilename);
    let command;

    if (LOCAL_BUILD === false) {
        command = [
            `source /etc/profile`,
            `sudo systemctl stop cudos-explorer.service`,
            `sudo systemctl stop cudos-faucet.service`,
            `mongo test --eval "db.dropDatabase()"`,
            `cd ${secrets.serverPath}`,
            `rm -R ./src`,
            `mkdir ./src`,
            `rm -R ./faucet`,
            `mkdir ./faucet`,
            `rm -Rf ./explorer`,
            `mkdir ./explorer`,
            // `tar -zxf ${filePath} -C ./src`,
            `unzip -q ${filePath} -d ./src`,
            `rm ${filePath}`,
            `cd ./src`,
            `cd ./project-faucet-cli`,
            `make`,
            `cp ./bin/cudos-noded "$GOPATH/bin"`,
            `sudo cp ./bin/libwasmvm.so "/usr/lib"`,
            `chmod +x "$GOPATH/bin/cudos-noded"`,
            `chmod +x ./run-testnet.sh`,
            `sed -i 's/\r$//' ./run-testnet.sh`,
            `cp ./run-testnet.sh ../../faucet`,
            `cd ../project-explorer`,
            `npm i`,
            `meteor build ../meteor-build/ --architecture os.linux.x86_64 --server-only`,
            `cd ../meteor-build`,
            `tar -zxf ./project-explorer.tar.gz`,
            `cd ./bundle/programs/server`,
            `npm i`,
            `cd ../../../`,
            `cp -r ./bundle/* ../../explorer`,
            `cd ../../explorer`,
            `cp ../src/project-explorer/run-testnet.sh ./`,
            `chmod +x ./run-testnet.sh`,
            `sed -i 's/\r$//' ./run-testnet.sh`,
            `sudo systemctl start cudos-explorer.service`,
            `sudo systemctl start cudos-faucet.service`,
            `cd ${secrets.serverPath}`,
            `rm -Rf ./src/*`
        ];
    } else {
        command = [
            `source /etc/profile`,
            `sudo systemctl stop cudos-explorer.service`,
            `sudo systemctl stop cudos-faucet.service`,
            `mongo test --eval "db.dropDatabase()"`,
            `cd ${secrets.serverPath}`,
            `rm -R ./src`,
            `mkdir ./src`,
            `rm -R ./faucet`,
            `mkdir ./faucet`,
            `rm -Rf ./explorer`,
            `mkdir ./explorer`,
            `unzip -q ${filePath} -d ./src`,
            `rm ${filePath}`,
            `cd ./src`,
            `cd ./project-faucet-cli`,
            `make`,
            `cp ./bin/cudos-noded "$GOPATH/bin"`,
            `sudo cp ./bin/libwasmvm.so "/usr/lib"`,
            `chmod +x "$GOPATH/bin/cudos-noded"`,
            `chmod +x ./run-testnet.sh`,
            `sed -i 's/\r$//' ./run-testnet.sh`,
            `cp ./run-testnet.sh ../../faucet`,
            `cd ../project-explorer`,
            `cd ./programs/server`,
            `npm i`,
            `cd ../../`,
            `cp -r ./* ../../explorer`,
            `cd ../../explorer`,
            `chmod +x ./run-testnet.sh`,
            `sed -i 's/\r$//' ./run-testnet.sh`,
            `sudo systemctl start cudos-explorer.service`,
            `sudo systemctl start cudos-faucet.service`,
            `cd ${secrets.serverPath}`,
            `rm -Rf ./src/*`
        ];
    }

    command = command.join(' && ');

    conn.on('ready', () => {
        console.log('Client :: ready');
        conn.exec(command, (err, stream) => {
            if (err) {
                throw err;
            }

            stream.on('close', (code, signal) => {
                console.log(`Stream :: close :: code: ${code}, signal: ${signal}`);
                conn.end();
            }).on('data', (data) => {
                console.log(`STDOUT: ${data}`);
            }).stderr.on('data', (data) => {
                console.log(`STDERR: ${data}`);
            });
        });
    });

    conn.connect({
        host: secrets.host,
        port: secrets.port,
        username: secrets.username,
        passphrase: secrets.keyPass,
        privateKey: (await asyncFs.readFile(secrets.privateKey)).toString(),
        path: secrets.serverPath,
    });
}

main();
