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

const OwnerManager = ({jwtToken}) => {
  const [teaOwners, setTeaOwners] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  const deleteOwner = id => {
    fetch(serverURL + '/owner/' + id, {
      method: 'DELETE',
      headers: {
        Token: jwtToken,
      },
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        } else {
          let newTeaOwners = [...teaOwners];
          newTeaOwners[0].data = newTeaOwners[0].data.filter(
            owner => owner.id !== id,
          );
          setTeaOwners(newTeaOwners);
          ToastAndroid.showWithGravityAndOffset(
            'Owner successfully deleted!',
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
          'Error deleting owner',
          'Please check if they own any teas!',
        );
      });
  };

  const addOwner = (name, _, textInput) => {
    if (!name) {
      Alert.alert('Error', 'Please enter a name for the new owner!');
      return;
    } else if (teaOwners[0].data.some(owner => owner.name === name)) {
      Alert.alert('Error', 'An owner with that name already exists!');
      return;
    }

    addOwnerAPI(name, textInput);
  };

  const addOwnerAPI = (name, textInput) => {
    fetch(serverURL + '/owner', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Token: jwtToken,
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
        let newTeaOwners = [...teaOwners];
        newTeaOwners[0].data.push({id: json.id, name});
        newTeaOwners[0].data.sort((a, b) =>
          a.name > b.name ? 1 : b.name > a.name ? -1 : 0,
        ); // Sort the owners
        setTeaOwners(newTeaOwners);
        ToastAndroid.showWithGravityAndOffset(
          'Owner successfully added!',
          ToastAndroid.SHORT,
          ToastAndroid.BOTTOM,
          0,
          150,
        );
        textInput.clear();
      })
      .catch(error => {
        console.warn(error);
        Alert.alert('Error adding owner', 'Please try again.');
      });
  };

  useEffect(() => {
    fetch(serverURL + '/owners', {
      headers: {
        Token: jwtToken,
      },
    })
      .then(response => response.json())
      .then(json => {
        let newTeaOwners = [{title: '', data: []}];
        json.forEach(owner => {
          newTeaOwners[0].data.push(owner);
        });
        newTeaOwners[0].data.sort((a, b) =>
          a.name > b.name ? 1 : b.name > a.name ? -1 : 0,
        );
        setTeaOwners(newTeaOwners);
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
          sections={teaOwners}
          renderItem={({item}) => (
            <ListItem item={item} deleteFunc={deleteOwner} />
          )}
          keyExtractor={(item, index) => item + index}
          renderSectionFooter={({section: {title}}) => (
            <AddSectionItem
              placeholderText={'Add Owner...'}
              addFunc={addOwner}
              sectionID={title.id}
            />
          )}
        />
      )}
    </View>
  );
};

export default OwnerManager;
