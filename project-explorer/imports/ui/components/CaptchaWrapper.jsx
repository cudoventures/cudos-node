/* global grecaptcha */

import React from 'react';
import PropTypes from 'prop-types';

import '../stylesheets/captcha-wrapper.css';

export default class CaptchaWrapper extends React.Component {

    constructor(props) {
        super(props);
        
        this.captchaId = null;
        this.executeCallback = null;

        this.nodes = {
            'captcha': React.createRef(),
        };

        this.onData = this.onData.bind(this);
        this.onError = this.onData.bind(this, null);
    }

    onData(response) {

        const callback = this.executeCallback;
        if (callback === null) {
            return;
        }

        this.executeCallback = null;

        if (response === null) {
            // Alert.show('Wrong catpcha, please reload the page and try again', () => {
            //     window.location.reload();
            // });
            console.error('Wrong catpcha, please reload the page and try again');
            return;
        }

        callback(response);
    }

    execute(callback) {

        grecaptcha.execute(this.captchaId);
        this.executeCallback = callback;
    }

    // getResponse() {
    //     return grecaptcha.getResponse(this.captchaId);
    // }

    reset() {
        grecaptcha.reset(this.captchaId);
    }

    componentDidMount() {
        const scriptN = document.createElement('script');
        scriptN.async = true;
        scriptN.defer = true;
        scriptN.src = 'https://www.google.com/recaptcha/api.js?render=explicit';
        scriptN.addEventListener('load', () => {
            if (grecaptcha === undefined) {
                Alert.show('Your browser is blocking Google\' recaptcha library. Please use different browser or enable Google\' recaptcha library', () => {
                    window.location.reload();
                });
                return;
            }

            grecaptcha.ready(() => {
                window.captcha = this;
                this.captchaId = grecaptcha.render(this.nodes.captcha.current, {
                    'sitekey': Meteor.settings.public.CAPTCHA_FRONTEND_KEY,
                    'size': 'invisible',
                    'theme': 'light',
                    'callback': this.onData,
                    'error-callback': this.onError,
                    // 'callback': (data) => {
                    //     console.log('1', data);
                    //     if (this.props.onCheck !== null) {
                    //         this.props.onCheck();
                    //     }
                    // },
                });
            });
        });
        document.head.appendChild(scriptN);
    }

    render() {
        return (
            <div ref = { this.nodes.captcha } className = { `CaptchaWrapper ${this.props.className}` } />
        )
    }

}

CaptchaWrapper.defaultProps = {
    'className': '',
    'onCheck': null,
};

CaptchaWrapper.propTypes = {
    'className': PropTypes.string,
    'onCheck': PropTypes.func,
};
