import qs from 'querystring';
import React,{ Component } from 'react';
import { HTTP } from 'meteor/http'
import {
    Badge,
    Button,
    Collapse,
    Navbar,
    NavbarToggler,
    NavbarBrand,
    Nav,
    NavItem,
    NavLink,
    // Input,
    // InputGroup,
    // InputGroupAddon,
    // Button,
    UncontrolledDropdown,
    UncontrolledPopover,
    PopoverBody,
    DropdownToggle,
    DropdownMenu,
    DropdownItem
} from 'reactstrap';
import { Link } from 'react-router-dom';
import SearchBar from './SearchBar.jsx';
import i18n from 'meteor/universe:i18n';
import LedgerModal from '../ledger/LedgerModal.jsx';
import Account from './Account.jsx';
import { DirectSecp256k1HdWallet, Registry } from "@cosmjs/proto-signing";
import { assertIsBroadcastTxSuccess, SigningStargateClient, StargateClient } from "@cosmjs/stargate";
import { MsgDelegate } from "@cosmjs/stargate/build/codec/cosmos/staking/v1beta1/tx"; 

const chainId = Meteor.settings.public.chainId;
const shainName = Meteor.settings.public.chainName;
const rpc = Meteor.settings.public.urls.rpc;
const rest = Meteor.settings.public.urls.rest;
const coinDenom = Meteor.settings.public.coins[0].denom;

const T = i18n.createComponent();

// Firefox does not support named group yet
// const SendPath = new RegExp('/account/(?<address>\\w+)/(?<action>send)')
// const DelegatePath = new RegExp('/validators?/(?<address>\\w+)/(?<action>delegate)')
// const WithdrawPath = new RegExp('/account/(?<action>withdraw)')

const SendPath = new RegExp('/account/(\\w+)/(send)')
const DelegatePath = new RegExp('/validators?/(\\w+)/(delegate)')
const WithdrawPath = new RegExp('/account/(withdraw)')

const getUser = () => localStorage.getItem(CURRENTUSERADDR)

export default class Header extends Component {
    constructor(props) {
        super(props);

        this.toggle = this.toggle.bind(this);

        this.state = {
            isOpen: false,
            networks: "",
            version: "-"
        };
    }

    toggle() {
        this.setState({
            isOpen: !this.state.isOpen
        }, ()=>{
            // console.log(this.state.isOpen);
        });
    }

    toggleSignIn = (value) => {
        this.setState(( prevState) => {
            return {isSignInOpen: value!=undefined?value:!prevState.isSignInOpen}
        })
    }

    handleLanguageSwitch(lang, e) {
        i18n.setLocale(lang)
    }

    componentDidMount(){
        const url = Meteor.settings.public.networks
        if (url){
            try{
                HTTP.get(url, null, (error, result) => {
                    if (result.statusCode == 200){
                        let networks = JSON.parse(result.content);
                        if (networks.length > 0){
                            this.setState({
                                networks: <DropdownMenu>{
                                    networks.map((network, i) => {
                                        return <span key={i}>
                                            <DropdownItem header><img src={network.logo} /> {network.name}</DropdownItem>
                                            {network.links.map((link, k) => {
                                                return <DropdownItem key={k} disabled={link.chain_id == chainId}>
                                                    <a href={link.url} target="_blank">{link.chain_id} <Badge size="xs" color="secondary">{link.name}</Badge></a>
                                                </DropdownItem>})}
                                            {(i < networks.length - 1)?<DropdownItem divider />:''}
                                        </span>

                                    })
                                }</DropdownMenu>
                            })
                        }
                    }
                })
            }
            catch(e){
                console.warn(e);
            }
        }

        Meteor.call('getVersion', (error, result) => {
            if (result) {
                this.setState({
                    version:result
                })
            }
        })
    }

    signOut = () => {
        localStorage.removeItem(CURRENTUSERADDR);
        localStorage.removeItem(CURRENTUSERPUBKEY);
        this.props.refreshApp();
    }

    shouldLogin = () => {
        let pathname = this.props.location.pathname
        let groups;
        let match = pathname.match(SendPath) || pathname.match(DelegatePath)|| pathname.match(WithdrawPath);
        if (match) {
            if (match[0] === '/account/withdraw') {
                groups = {action: 'withdraw'}
            } else {
                groups = {address: match[1], action: match[2]}
            }
        }
        let params = qs.parse(this.props.location.search.substr(1))
        return groups || params.signin != undefined
    }

    handleLoginConfirmed = (success) => {
        let groups = this.shouldLogin()
        if (!groups) return
        let redirectUrl;
        let params;
        if (groups) {
            let { action, address } = groups;
            params = {action}
            switch (groups.action) {
            case 'send':
                params.transferTarget = address
                redirectUrl = `/account/${address}`
                break
            case 'withdraw':
                redirectUrl = `/account/${getUser()}`
                break;
            case 'delegate':
                redirectUrl = `/validators/${address}`
                break;
            }
        } else {
            let location = this.props.location;
            params = qs.parse(location.search.substr(1))
            redirectUrl = params.redirect?params.redirect:location.pathname;
            delete params['redirectUrl']
            delete params['signin']
        }

        let query = success?`?${qs.stringify(params)}`:'';
        this.props.history.push(redirectUrl + query)
    }

    connectKeplr = async () => {
        if (!window.getOfflineSigner || !window.keplr) {
            alert("Please install keplr extension");
        } else {
            if (window.keplr.experimentalSuggestChain) {
                try {
                    await window.keplr.experimentalSuggestChain({
                        // Chain-id of the Cosmos SDK chain.
                        chainId: chainId,
                        // The name of the chain to be displayed to the user.
                        chainName: shainName,
                        // RPC endpoint of the chain.
                        rpc: rpc,
                        // REST endpoint of the chain.
                        rest: rest,
                        // Staking coin information
                        stakeCurrency: {
                            // Coin denomination to be displayed to the user.
                            coinDenom: coinDenom,
                            // Actual denom (i.e. uatom, uscrt) used by the blockchain.
                            coinMinimalDenom: coinDenom,
                            // # of decimal points to convert minimal denomination to user-facing denomination.
                            coinDecimals: 6,
                            // (Optional) Keplr can show the fiat value of the coin if a coingecko id is provided.
                            // You can get id from https://api.coingecko.com/api/v3/coins/list if it is listed.
                            // coinGeckoId: ""
                        },
                        // (Optional) If you have a wallet webpage used to stake the coin then provide the url to the website in `walletUrlForStaking`.
                        // The 'stake' button in Keplr extension will link to the webpage.
                        walletUrlForStaking: "http://localhost:26657",
                        // The BIP44 path.
                        bip44: {
                            // You can only set the coin type of BIP44.
                            // 'Purpose' is fixed to 44.
                            coinType: 118,
                        },
                        bech32Config: {
                            bech32PrefixAccAddr: Meteor.settings.public.bech32PrefixAccAddr,
                            bech32PrefixAccPub: Meteor.settings.public.bech32PrefixAccPub,
                            bech32PrefixValAddr: Meteor.settings.public.bech32PrefixValAddr,
                            bech32PrefixValPub: Meteor.settings.public.bech32PrefixValPub,
                            bech32PrefixConsAddr: Meteor.settings.public.bech32PrefixConsAddr,
                            bech32PrefixConsPub: Meteor.settings.public.bech32PrefixConsPub
                        },
                        // List of all coin/tokens used in this chain.
                        currencies: [{
                            // Coin denomination to be displayed to the user.
                            coinDenom: coinDenom,
                            // Actual denom (i.e. uatom, uscrt) used by the blockchain.
                            coinMinimalDenom: coinDenom,
                            // # of decimal points to convert minimal denomination to user-facing denomination.
                            coinDecimals: 6,
                            // (Optional) Keplr can show the fiat value of the coin if a coingecko id is provided.
                            // You can get id from https://api.coingecko.com/api/v3/coins/list if it is listed.
                            // coinGeckoId: ""
                        }],
                        // List of coin/tokens used as a fee token in this chain.
                        feeCurrencies: [{
                            // Coin denomination to be displayed to the user.
                            coinDenom: coinDenom,
                            // Actual denom (i.e. uatom, uscrt) used by the blockchain.
                            coinMinimalDenom: coinDenom,
                            // # of decimal points to convert minimal denomination to user-facing denomination.
                            coinDecimals: 6,
                            // (Optional) Keplr can show the fiat value of the coin if a coingecko id is provided.
                            // You can get id from https://api.coingecko.com/api/v3/coins/list if it is listed.
                            // coinGeckoId: ""
                        }],
                        // (Optional) The number of the coin type.
                        // This field is only used to fetch the address from ENS.
                        // Ideally, it is recommended to be the same with BIP44 path's coin type.
                        // However, some early chains may choose to use the Cosmos Hub BIP44 path of '118'.
                        // So, this is separated to support such chains.
                        coinType: 118,
                        // (Optional) This is used to set the fee of the transaction.
                        // If this field is not provided, Keplr extension will set the default gas price as (low: 0.01, average: 0.025, high: 0.04).
                        // Currently, Keplr doesn't support dynamic calculation of the gas prices based on on-chain data.
                        // Make sure that the gas prices are higher than the minimum gas prices accepted by chain validators and RPC/REST endpoint.
                        gasPriceStep: {
                            low: 0.01,
                            average: 0.025,
                            high: 0.04
                        }
                    });
                } catch {
                    alert("Failed to suggest the chain");
                }
            } else {
                alert("Please use the recent version of keplr extension");
            }
        }
        // You should request Keplr to enable the wallet.
        // This method will ask the user whether or not to allow access if they haven't visited this website.
        // Also, it will request user to unlock the wallet if the wallet is locked.
        // If you don't request enabling before usage, there is no guarantee that other methods will work.
        await window.keplr.enable(chainId);

        const offlineSigner = window.getOfflineSigner(chainId);
        const account = (await offlineSigner.getAccounts())[0];
    
        localStorage.setItem(CURRENTUSERADDR, account.address);
        localStorage.setItem(CURRENTUSERPUBKEY, account.pubkey);
        this.props.refreshApp();
    }

    render() {
        let signedInAddress = getUser();
        return (
            <Navbar color="primary" dark expand="lg" fixed="top" id="header">
                <NavbarBrand tag={Link} to="/"><img src="/img/big-dipper-icon-light.svg" className="img-fluid logo"/> <span className="d-none d-xl-inline-block"><T>navbar.siteName</T>&nbsp;</span><Badge color="secondary">{this.state.version}</Badge> </NavbarBrand>
                <UncontrolledDropdown className="d-inline text-nowrap">
                    <DropdownToggle caret={(this.state.networks !== "")} tag="span" size="sm" id="network-nav">{Meteor.settings.public.chainId}</DropdownToggle>
                    {this.state.networks}
                </UncontrolledDropdown>
                <SearchBar id="header-search" history={this.props.history} />
                <NavbarToggler onClick={this.toggle} />
                <Collapse isOpen={this.state.isOpen} navbar>
                    <Nav className="ml-auto text-nowrap" navbar>
                        <NavItem>
                            <NavLink tag={Link} to="/validators"><T>navbar.validators</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/blocks"><T>navbar.blocks</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/transactions"><T>navbar.transactions</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/proposals"><T>navbar.proposals</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/voting-power-distribution"><T>navbar.votingPower</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/faucet"><T>navbar.faucet</T></NavLink>
                        </NavItem>
                        <NavItem id="user-acconut-icon">
                            {!signedInAddress?<Button className="sign-in-btn" color="link" size="lg" onClick={() => this.connectKeplr()}><i className="material-icons">vpn_key</i></Button>:
                                <span>
                                    <span className="d-lg-none">
                                        <i className="material-icons large d-inline">account_circle</i>
                                        <Link to={`/account/${signedInAddress}`}> {signedInAddress}</Link>
                                        <Button className="float-right" color="link" size="sm" onClick={this.signOut}><i className="material-icons">exit_to_app</i></Button>
                                    </span>
                                    <span className="d-none d-lg-block">
                                        <i className="material-icons large">account_circle</i>
                                        <UncontrolledPopover className="d-none d-lg-block" trigger="legacy" placement="bottom" target="user-acconut-icon">
                                            <PopoverBody>
                                                <div className="text-center"> 
                                                    <p><T>accounts.signInText</T></p>
                                                    <p><Link className="text-nowrap" to={`/account/${signedInAddress}`}>{signedInAddress}</Link></p>
                                                    <Button className="float-right" color="link" onClick={this.signOut}><i className="material-icons">exit_to_app</i><span> <T>accounts.signOut</T></span></Button>
                                                </div>
                                            </PopoverBody>
                                        </UncontrolledPopover>
                                    </span>
                                </span>}
                            <LedgerModal isOpen={this.state.isSignInOpen} toggle={this.toggleSignIn} refreshApp={this.props.refreshApp} handleLoginConfirmed={this.shouldLogin()?this.handleLoginConfirmed:null}/>
                        </NavItem>
                        {/* <NavItem>
                            <UncontrolledDropdown inNavbar>
                                <DropdownToggle nav caret>
                                    <T>navbar.lang</T>
                                </DropdownToggle>
                                <DropdownMenu right>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('en-US', e)}><T>navbar.english</T></DropdownItem>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('es-ES', e)}><T>navbar.spanish</T></DropdownItem>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('it-IT', e)}><T>navbar.italian</T></DropdownItem>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('pl-PL', e)}><T>navbar.polish</T></DropdownItem>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('ru-RU', e)}><T>navbar.russian</T></DropdownItem>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('zh-Hant', e)}><T>navbar.chinese</T></DropdownItem>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('zh-Hans', e)}><T>navbar.simChinese</T></DropdownItem>
                                </DropdownMenu>
                            </UncontrolledDropdown>
                        </NavItem> */}
                    </Nav>
                </Collapse>
            </Navbar>
        );
    }
}