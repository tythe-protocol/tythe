import PropTypes from 'prop-types';
import React, { Component } from 'react';

export default class List extends Component {
    constructor(props) {
        super(props);
        this.state = {
            deps: null,
        };
    }

    componentDidMount() {
        this.controller = new AbortController();
        fetch("/-/list?r=" + this.props.repo, {
           signal: this.controller.signal, 
        }).then(r => r.json()).then(deps => {
            this.setState({
                deps: deps,
            });
        }).catch(err => console.error);
    }

    componentWillUnmount() {
        this.controller.abort();
    }

    render() {
        return <div>
            <div style={{fontSize: "1.5em", marginBottom: "1em"}}>{this.props.repo}</div>
            {table(this.state.deps)}
        </div>;
    }
}

List.propTypes = {
    repo: PropTypes.string.isRequired,
};

function table(deps) {
    if (!deps) {
        return "Loading...";
    }
    return <table cellSpacing="0" style={{borderCollapse: "collapse"}}>
        <tbody>
            <tr key="header">
                {th("Type")}
                {th("Name")}
                {th("Address")}
            </tr>
            {deps.map(tr)}
        </tbody>
    </table>
}

function tr(dep) {
    return <tr key={key(dep)}>
        {td(type(dep.type))}
        {td(dep.name)}
        {td(config(dep.config))}
    </tr>
}

function th(child) {
    return <th style={{border: "1px solid black", padding: "0.25em"}}>{child}</th>;
}

function td(child) {
    return <td style={{border: "1px solid black", padding: "0.25em"}}>{child}</td>;
}

function type(t) {
    switch (t) {
        case 1:
            return "Go"
        case 2:
            return "NPM"
        default:
            return "??"
    }
}

function config(c) {
    if (!c) {
        return "<none>";
    } else {
        return JSON.stringify(c);
    }
}

function key(dep) {
    return dep.type + ":" + dep.name;
}
