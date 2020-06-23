import React, { useState, useEffect } from 'react';
import { serverURL } from '../app.json';
import {
  View,
  Text,
  ActivityIndicator,
} from 'react-native';
import Button from 'react-native-button';

const RandomSelector = (props) => {

  const [teas, setTeas] = useState([]);
  const [selectedTea, setSelectedTea] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  
  const selectRandomTea = () => {
    setSelectedTea(Math.floor(Math.random()*teas.length));
  };

  useEffect(() => {
    fetch(serverURL + '/teas')
    .then((response) => response.json())
    .then((json) => {
      setTeas(json);
    })
    .then(() => selectRandomTea())
    .catch((error) => console.error(error))
    .finally(() => {
      setIsLoading(false);
    })
  }, []);

  return (
      <View style={props.styles.sectionContainer}>
        <Text style={props.styles.sectionTitle}>The selected tea is:</Text>
        {isLoading || selectedTea === null ? <ActivityIndicator/> : (
          <Text style={props.styles.sectionDescription}>{ teas[selectedTea].name }</Text>
        )}
        <Button
          onPress={selectRandomTea}
          containerStyle={ props.styles.buttonContainer }
          style={props.styles.button }>
          Select Another Tea
        </Button>
      </View>
  );
};

export default RandomSelector;
