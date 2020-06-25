import React from 'react';
import {View, Text, StyleSheet} from 'react-native';
import CheckBox from '@react-native-community/checkbox';

const CheckBoxGroup = ({items, states, updateFunc}) => {
  return (
    <View style={styles.container}>
      {items.map((value, index) => {
        return (
          <View style={styles.checkboxRow} key={index}>
            <CheckBox
              value={states[index]}
              onValueChange={() => updateFunc(index)}
            />
            <Text style={styles.label}>{value.name}</Text>
          </View>
        );
      })}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    justifyContent: 'center',
    alignItems: 'center',
    paddingTop: 10,
  },
  checkboxRow: {
    flexDirection: 'row',
    paddingBottom: 5,
  },
  label: {
    fontSize: 18,
    textAlignVertical: 'center',
  },
});

export default CheckBoxGroup;
