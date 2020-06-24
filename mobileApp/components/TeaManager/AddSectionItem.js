import React, {useState} from 'react';
import {TextInput, View, StyleSheet, Dimensions} from 'react-native';
import Button from 'react-native-button';
import Icon from 'react-native-vector-icons/dist/FontAwesome5';

const AddSectionItem = ({placeholderText, addFunc, sectionID}) => {
  const [addText, setAddText] = useState('');

  const onChange = textValue => setAddText(textValue);

  return (
    <View style={styles.inputView}>
      <TextInput
        placeholder={placeholderText}
        style={styles.input}
        onChangeText={onChange}
        selection="center"
        ref={input => {
          this.textInput = input;
        }}
      />
      <Button
        containerStyle={styles.btnContainer}
        onPress={() => {
          addFunc(addText, sectionID, this.textInput);
        }}>
        <Icon name="plus" size={20} color="white" style={styles.icon} />
      </Button>
    </View>
  );
};

const styles = StyleSheet.create({
  inputView: {
    backgroundColor: 'white',
    flexDirection: 'row',
    paddingTop: 10,
    paddingBottom: 10,
  },
  input: {
    height: 40,
    padding: 5,
    marginLeft: 5,
    fontSize: 18,
    borderTopLeftRadius: 50,
    borderBottomLeftRadius: 50,
    width: Math.round(Dimensions.get('window').width) - 50,
    backgroundColor: 'lightgray',
    textAlign: 'center',
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

export default AddSectionItem;
