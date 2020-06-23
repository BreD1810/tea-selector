import React, {useState, useEffect} from 'react';
import {
  View,
  FlatList,
  ActivityIndicator,
  Alert,
  ToastAndroid,
} from 'react-native';
import TeaListItem from './TeaListItem';
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
    </View>
  );
};

export default TeaManager;
