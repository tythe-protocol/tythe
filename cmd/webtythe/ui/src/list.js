import PropTypes from 'prop-types';
import React, { Component } from 'react';
import clarinet from 'clarinet';

export default class List extends Component {
    constructor(props) {
        super(props);
        this.state = {
            deps: null,
        };
    }

    componentDidMount() {
        const parser = clarinet.parser();

        const stack = [];
        let key;
        const peek = () => {
            if (stack.length == 0) {
                return null;
            }
            return stack[stack.length - 1];
        }
        const isArray = v => v.constructor == Array;
        const append = v => {
            const p = peek();
            if (p) {
                if (isArray(p)) {
                    p.push(v);
                } else {
                    p[key] = v;
                }
            }
            return v;
        }

        parser.onvalue = v => {
            append(v);
        };
        parser.onopenobject = k => {
            stack.push(append({}));
            key = k;
        };
        parser.onkey = k => {
            key = k;
        };
        parser.oncloseobject = () => {
            const o = stack.pop();
            if (stack.length == 1) {
                console.log(o);
            }
        };
        parser.onopenarray = () => {
            stack.push(append([]));
        };
        parser.onclosearray = () => {
            stack.pop();
        }

        this.controller = new AbortController();
        fetch("/-/list?r=" + this.props.repo, {
           signal: this.controller.signal, 
        })
        .then(r => {
            const decoder = new TextDecoder();
            const reader = r.body.getReader();
            let chunk = ({done, value}) => {
                if (!done) {
                    const text = decoder.decode(value, {stream:true});
                    parser.write(text);
                    reader.read().then(chunk);
                }
            };
            reader.read().then(chunk);
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
