import React from 'react';
import {
  SafeAreaView,
  StyleSheet,
  View,
} from 'react-native';
import RandomSelector from './components/RandomSelector';

const App: () => React$Node = () => {
  return (
    <>
      <SafeAreaView>
          <View style={styles.body}>
            <RandomSelector styles={ styles } />
          </View>
      </SafeAreaView>
    </>
  );
};

const styles = StyleSheet.create({
  sectionContainer: {
    marginTop: 32,
    paddingHorizontal: 24,
  },
  sectionTitle: {
    fontSize: 32,
    fontWeight: '600',
    color: 'black',
    alignSelf: 'center'
  },
  sectionDescription: {
    marginTop: 8,
    fontSize: 24,
    fontWeight: '400',
    color: 'black',
    alignSelf: 'center'
  },
  button: {
    fontSize: 20,
    color: 'white',
  },
  buttonContainer: {
    padding: 15,
    marginTop: 20,
    height: 60,
    width: 200,
    overflow: 'hidden',
    borderRadius: 4,
    backgroundColor: 'dodgerblue',
    alignSelf: 'center'
  },
});

export default App;
