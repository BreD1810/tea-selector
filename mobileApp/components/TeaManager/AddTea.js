import React, {useState} from 'react';
import {TextInput, View, StyleSheet} from 'react-native';
import Button from 'react-native-button';
import Icon from 'react-native-vector-icons/dist/FontAwesome5';

const AddTea = ({addTea}) => {
  const [addText, setAddText] = useState('');

  const onChange = textValue => setAddText(textValue);

  return (
    <View>
      <TextInput
        placeholder="Add Tea..."
        style={styles.input}
        onChangeText={onChange}
      />
      <Button
        style={styles.btn}
        containerStyle={styles.btnContainer}
        onPress={() => addTea(addText)}>
        <Icon name="plus" size={20} color="white" />
        Add Tea
      </Button>
    </View>
  );
};

const styles = StyleSheet.create({
  input: {
    height: 60,
    padding: 8,
    fontSize: 16,
    backgroundColor: 'lightgray',
  },
  btn: {
    fontSize: 20,
    color: 'white',
    paddingLeft: 10,
  },
  btnContainer: {
    alignSelf: 'center',
    flexDirection: 'row',
    justifyContent: 'center',
    padding: 15,
    marginTop: 5,
    marginBottom: 15,
    height: 60,
    width: 200,
    overflow: 'hidden',
    borderRadius: 4,
    backgroundColor: 'dodgerblue',
  },
});

export default AddTea;
