import React, { useState, useEffect } from 'react';
import { serverURL } from '../app.json';
import {
  View,
  Text,
  StyleSheet,
  ActivityIndicator,
} from 'react-native';
import Button from 'react-native-button';

const HomePage = (props) => {

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
      <View style={ styles.container }>
        <Text style={ styles.title }>The selected tea is:</Text>
        {isLoading || selectedTea === null ? <ActivityIndicator/> : (
          <Text style={ styles.tea }>{ teas[selectedTea].name }</Text>
        )}
        <Button
          onPress={ selectRandomTea }
          containerStyle={ styles.btnContainer }
          style={ styles.btn }>
          Select Another Tea
        </Button>
      </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
  },
  title: {
    fontSize: 32,
    fontWeight: '600',
    color: 'black',
  },
  tea: {
    marginTop: 8,
    fontSize: 24,
    fontWeight: '400',
    color: 'black',
  },
  btn: {
    fontSize: 20,
    color: 'white',
  },
  btnContainer: {
    padding: 15,
    marginTop: 20,
    height: 60,
    width: 200,
    overflow: 'hidden',
    borderRadius: 4,
    backgroundColor: 'dodgerblue',
  },
});

export default HomePage;
