import React from 'react';
import {Text, View, StyleSheet} from 'react-native';
import Button from 'react-native-button';

const ManageButtons = ({navigation}) => {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Manage</Text>
      <Button
        style={styles.btn}
        containerStyle={styles.btnContainer}
        onPress={() => navigation.navigate('TeaManager')}>
        Manage Teas
      </Button>
      <Button
        style={styles.btn}
        containerStyle={styles.btnContainer}
        onPress={() => navigation.navigate('Temp')}>
        Manage Tea Types
      </Button>
      <Button
        style={styles.btn}
        containerStyle={styles.btnContainer}
        onPress={() => navigation.navigate('Temp')}>
        Manage Owners
      </Button>
      <Button
        style={styles.btnSmall}
        containerStyle={styles.btnContainerSmall}
        onPress={() => navigation.navigate('Temp')}>
        Manage Ownership
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
  btn: {
    fontSize: 20,
    color: 'white',
  },
  btnSmall: {
    fontSize: 16,
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
  btnContainerSmall: {
    padding: 18,
    marginTop: 20,
    height: 60,
    width: 200,
    overflow: 'hidden',
    borderRadius: 4,
    backgroundColor: 'dodgerblue',
  },
});

export default ManageButtons;
