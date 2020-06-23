import React, {useState, useEffect} from 'react';
import {
  View,
  FlatList,
  ActivityIndicator
} from 'react-native';
import TeaListItem from './TeaListItem';
import { serverURL } from '../../app.json';

const TeaManager = () => {

  const [teas, setTeas] = useState([
    {id: '1', name: 'Test'},
    {id: '2', name: 'Test 2'},
    {id: '3', name: 'Test 3'},
  ]);
  const [isLoading, setIsLoading] = useState(true);
  
  useEffect(() => {
    fetch(serverURL + '/teas')
    .then((response) => response.json())
    .then((json) => {
      setTeas(json);
    })
    .catch((error) => console.error(error))
    .finally(() => {
      setIsLoading(false);
    })
  }, []);
  
  return (
    <View>
      {isLoading ? <ActivityIndicator/> : (
        <FlatList 
          data={teas}
          renderItem={({item}) => <TeaListItem item={item} />}
          keyExtractor={item => item.id}
        />
      )}
    </View>
  );
};

export default TeaManager;
