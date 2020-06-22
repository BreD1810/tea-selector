import React, { Component } from 'react';
import { serverURL } from '../app.json';
import {
  View,
  Text,
  ActivityIndicator,
} from 'react-native';
import Button from 'react-native-button';

export default class RandomSelector extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      teas: [],
      selectedTea: null,
      isLoading: true
    };
  }

  componentDidMount() {
    fetch(serverURL + '/teas')
      .then((response) => response.json())
      .then((json) => {
        this.setState({ teas: json });
      })
      .then(() => this.selectRandomTea())
      .catch((error) => console.error(error))
      .finally(() => {
        this.setState({ isLoading: false });
      })
  }

  selectRandomTea() {
    const teas = this.state.teas;
    this.setState( {selectedTea: teas[Math.floor(Math.random()*teas.length)] });
  }

  render() {
      return (
          <View style={this.props.styles.sectionContainer}>
            <Text style={this.props.styles.sectionTitle}>The selected tea is:</Text>
            {this.state.isLoading || this.state.selectedTea === null ? <ActivityIndicator/> : (
              <Text style={this.props.styles.sectionDescription}>{ this.state.selectedTea.name }</Text>
            )}
            <Button
              onPress={() => this.selectRandomTea()}
              containerStyle={ this.props.styles.buttonContainer }
              style={this.props.styles.button }>
              Select Another Tea
            </Button>
          </View>
      );
  }
}
