import React, { Component } from 'react';

export default class Home extends Component {
    render() {
        return (
            <div>
                <div style={{
                    margin: "1em 6em 1em",
                    fontSize: "1em",
                    color: "#444",
                }}>Enter the URL to your Git repository:</div>
                <form action="list" method="GET" style={{
                    margin: "1em 6em",
                }}>
                    <input style={{
                    fontSize: "2em",
                    padding: "0.25em",
                    width: "100%",
                    boxSizing: "border-box",
                    }} type="text" name="r" defaultValue="https://github.com/tythe-protocol/tythe"/>
                    <input style={{
                    fontSize: "1em",
                    margin: "1em 0",
                    }} type="submit" value="Start!"/>
                </form>
            </div>
        );
    }
}
