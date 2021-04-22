class ArgvHelper {

    static parse() {
        for (var value, split, i = 2;  i < process.argv.length;  ++i) {
            value = process.argv[i];
            if (value.indexOf('-target') === 0) {
                split = value.split('=');
                if (split.length === 2)
                    ArgvHelper.TARGET = split[1].trim();

                continue;
            }
        }

        switch (ArgvHelper.TARGET) {
            case 'staging':
                break;
            case 'production':
                break;
        }
    }

}

ArgvHelper.TARGET = null;

ArgvHelper.parse();

module.exports = ArgvHelper;
