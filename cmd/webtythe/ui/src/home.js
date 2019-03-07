import React, { Component } from 'react';

export default class Home extends Component {
    render() {
        return (
            <div>
                <div>Enter the URL to your Git repository:</div>
                <form method="GET" style={{
                    margin: "1em 0",
                }}>
                    <input style={{
                    fontSize: "2em",
                    padding: "0.25em",
                    width: "100%",
                    boxSizing: "border-box",
                    }} type="text" name="list" defaultValue="https://github.com/tythe-protocol/tythe"/>
                    <input style={{
                    fontSize: "1em",
                    margin: "1em 0",
                    }} type="submit" value="Start!"/>
                </form>
            </div>
        );
    }
}
