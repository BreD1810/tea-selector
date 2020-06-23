import React from 'react';
import {Text, TouchableOpacity, View, StyleSheet} from 'react-native';
import Icon from 'react-native-vector-icons/dist/FontAwesome5';

const TeaListItem = ({item, deleteTea}) => {
  return (
    <TouchableOpacity style={styles.listItem}>
      <View style={styles.listItemView}>
        <Text style={styles.listItemText}>{item.name}</Text>
        <Icon
          name="trash"
          size={20}
          color="firebrick"
          onPress={() => deleteTea(item.id)}
        />
      </View>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  listItem: {
    padding: 15,
    backgroundColor: '#f8f8f8',
    borderWidth: 1,
    borderColor: '#eee',
  },
  listItemView: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  listItemText: {
    fontSize: 18,
  },
});

export default TeaListItem;
