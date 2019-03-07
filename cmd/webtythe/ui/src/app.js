import React, { Component } from 'react';
import Home from './home.js';

export default class App extends Component {
  render() {
    return (
      <div id="app">
        <header style={{
          padding: "3em 3em 2em",
          fontSize: " 2em",
          fontWeight: 800,
          color: "#444",
        }}>tythe<span style={{color:"#aaa"}}>.dev</span></header>
        <Home/>
      </div>
    );
  }
}
