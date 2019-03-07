import PropTypes from 'prop-types';
import React, { Component } from 'react';

import Home from './home.js';
import List from './list.js';

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
        <div style={{margin: "1em 6em 1em", color: "#444"}}>
          {content(this.props.url)}
        </div>
      </div>
    );
  }
}

App.propTypes = {
  url: PropTypes.instanceOf(URL).isRequired,
};

function content(url) {
  const list = url.searchParams.get("list");
  if (list) {
    return <List repo={list}/>
  }
  return <Home/>;
}
