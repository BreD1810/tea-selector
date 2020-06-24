import React, {useState, useEffect} from 'react';
import {
  View,
  SectionList,
  ActivityIndicator,
  Alert,
  ToastAndroid,
} from 'react-native';
import ListItem from './Lists/ListItem';
import AddSectionItem from './Lists/AddSectionItem';
import {serverURL} from '../../app.json';

const TeaTypeManager = () => {
  const [teaTypes, setTeaTypes] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  const deleteType = id => {
    fetch(serverURL + '/type/' + id, {
      method: 'DELETE',
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        } else {
          let newTeaTypes = [...teaTypes];
          newTeaTypes[0].data = newTeaTypes[0].data.filter(
            teaType => teaType.id !== id,
          );
          setTeaTypes(newTeaTypes);
          ToastAndroid.showWithGravityAndOffset(
            'Tea type successfully deleted!',
            ToastAndroid.SHORT,
            ToastAndroid.BOTTOM,
            0,
            150,
          );
        }
      })
      .catch(error => {
        console.warn(error);
        Alert.alert(
          'Error deleting tea type',
          'Please check if there are any teas of this type!',
        );
      });
  };

  const addTeaType = (name, _, textInput) => {
    if (!name) {
      Alert.alert('Error', 'Please enter a name for the new type!');
      return;
    } else if (teaTypes[0].data.some(teaType => teaType.name === name)) {
      Alert.alert('Error', 'That type of tea already exists!');
      return;
    }

    addTeaTypeAPI(name, textInput);
  };

  const addTeaTypeAPI = (name, textInput) => {
    fetch(serverURL + '/type', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({name: name}),
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        }
        return response.json();
      })
      .then(json => {
        let newTeaTypes = [...teaTypes];
        newTeaTypes[0].data.push({id: json.id, name});
        setTeaTypes(newTeaTypes);
        ToastAndroid.showWithGravityAndOffset(
          'Tea type successfully added!',
          ToastAndroid.SHORT,
          ToastAndroid.BOTTOM,
          0,
          150,
        );
        textInput.clear();
      })
      .catch(error => {
        console.warn(error);
        Alert.alert('Error adding type', 'Please try again.');
      });
  };

  useEffect(() => {
    fetch(serverURL + '/types')
      .then(response => response.json())
      .then(json => {
        let newTeaTypes = [{title: '', data: []}];
        json.forEach(teaType => {
          newTeaTypes[0].data.push(teaType);
        });
        setTeaTypes(newTeaTypes);
      })
      .catch(error => console.error(error))
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  return (
    <View>
      {isLoading ? (
        <ActivityIndicator />
      ) : (
        <SectionList
          sections={teaTypes}
          renderItem={({item}) => (
            <ListItem item={item} deleteFunc={deleteType} />
          )}
          keyExtractor={(item, index) => item + index}
          renderSectionFooter={({section: {title}}) => (
            <AddSectionItem
              placeholderText={'Add Type...'}
              addFunc={addTeaType}
              sectionID={title.id}
            />
          )}
        />
      )}
    </View>
  );
};

export default TeaTypeManager;
