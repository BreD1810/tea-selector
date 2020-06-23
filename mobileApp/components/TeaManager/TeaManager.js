import React, {useState, useEffect} from 'react';
import {
  View,
  FlatList,
  ActivityIndicator,
  Alert,
  ToastAndroid,
} from 'react-native';
import TeaListItem from './TeaListItem';
import AddTea from './AddTea';
import {serverURL} from '../../app.json';

const TeaManager = () => {
  const [teas, setTeas] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  const deleteTea = id => {
    fetch(serverURL + '/tea/' + id, {
      method: 'DELETE',
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        } else {
          setTeas(prevTeas => {
            return prevTeas.filter(tea => tea.id !== id);
          });
          ToastAndroid.showWithGravityAndOffset(
            'Tea successfully deleted!',
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
          'Error deleting tea',
          'Please check if someone owns this tea!',
        );
      });
  };

  const addTea = name => {
    if (!name) {
      Alert.alert('Error', 'Please enter a name for the new tea!');
      return;
    } else if (teas.some(tea => tea.name === name)){
      Alert.alert('Error', 'That tea already exists!');
      return;
    }

    addTeaAPI(name);
  };

  const addTeaAPI = name => {
    fetch(serverURL + '/tea', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({name: name, type: {id: 1}}), // TODO: Add in type selected by the user.
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        }
        return response.json();
      })
      .then(json => {
        setTeas(prevTeas => {
          return [...prevTeas, {id: json.id, name}];
        });
        ToastAndroid.showWithGravityAndOffset(
          'Tea successfully added!',
          ToastAndroid.SHORT,
          ToastAndroid.BOTTOM,
          0,
          150,
        );
      })
      .catch(error => {
        console.warn(error);
        Alert.alert('Error adding tea', 'Please try again.');
      });
  };

  useEffect(() => {
    fetch(serverURL + '/teas')
      .then(response => response.json())
      .then(json => {
        setTeas(json);
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
        <FlatList
          data={teas}
          renderItem={({item}) => (
            <TeaListItem item={item} deleteTea={deleteTea} />
          )}
          keyExtractor={item => item.id.toString()}
        />
      )}
      <AddTea addTea={addTea} />
    </View>
  );
};

export default TeaManager;
