import React, { Component } from 'react';

import {
  View,
  Text,
  ActivityIndicator
} from 'react-native';

export default class RandomSelector extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      teas: [],
      isLoading: true
    };
    this.getAllTeas();
  }

  getAllTeas() { 
    // fetch('https://192.168.1.184:7344/teas', {
    //   method: 'GET',
    // })
    console.log('FETCH')
    fetch('http://192.168.1.184:7344/teas')
      .then((response) => response.json())
      .then((json) => {
        this.setState({ teas: json });
      })
      .catch((error) => console.error(error))
      .finally(() => {
        this.setState({ isLoading: false });
      })
  }

  selectRandomTea() {
    const teas = this.state.teas;
    const selectedTea = teas[Math.floor(Math.random()*teas.length)];
    return <Text style={this.props.styles.sectionDescription}>{ selectedTea.name }</Text>
  }

  render() {
      return (
          <View style={this.props.styles.sectionContainer}>
            <Text style={this.props.styles.sectionTitle}>The selected tea is:</Text>
            {this.state.isLoading ? <ActivityIndicator/> : (
              this.selectRandomTea()
            )}
          </View>
      );
  }
}
