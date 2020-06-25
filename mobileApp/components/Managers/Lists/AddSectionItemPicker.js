import React, {useState} from 'react';
import {View, StyleSheet, Dimensions} from 'react-native';
import Button from 'react-native-button';
import RNPickerSelect from 'react-native-picker-select';
import Icon from 'react-native-vector-icons/dist/FontAwesome5';

const AddSectionItemPicker = ({promptText, inputItems, addFunc, sectionID}) => {
  const [selectedItem, setSelectedItem] = useState(null);
  const [items, setItems] = useState(inputItems);

  const onChange = value => setSelectedItem(value);

  const reset = () => setSelectedItem(null);

  return (
    <View style={styles.backgroundView}>
      <View style={styles.inputView}>
        <RNPickerSelect
          placeholder={{
            label: promptText,
            value: null,
            color: 'gray',
          }}
          items={items}
          value={selectedItem}
          onValueChange={onChange}
          useNativeAndroidPickerStyle={false}
          style={{
            ...styles.input,
            placeholder: {
              color: 'gray',
              paddingHorizontal: 5,
            },
          }}
        />
      </View>
      <Button
        containerStyle={styles.btnContainer}
        onPress={() => {
          addFunc(selectedItem, sectionID, reset);
        }}>
        <Icon name="plus" size={20} color="white" style={styles.icon} />
      </Button>
    </View>
  );
};

const styles = StyleSheet.create({
  backgroundView: {
    flexDirection: 'row',
    paddingTop: 10,
    paddingBottom: 10,
  },
  input: {
    width: Math.round(Dimensions.get('window').width) - 50,
  },
  inputView: {
    backgroundColor: 'lightgray',
    height: 40,
    width: Math.round(Dimensions.get('window').width) - 50,
    borderTopLeftRadius: 25,
    borderBottomLeftRadius: 25,
    paddingLeft: 15,
  },
  icon: {
    alignSelf: 'center',
    color: 'gray',
  },
  btnContainer: {
    paddingRight: 10,
    paddingTop: 10,
    height: 40,
    width: 40,
    borderTopRightRadius: 25,
    borderBottomRightRadius: 25,
    backgroundColor: 'lightgray',
  },
});

export default AddSectionItemPicker;
