import React, {useState, useEffect} from 'react';
import {
  View,
  Text,
  SectionList,
  ActivityIndicator,
  Alert,
  ToastAndroid,
  StyleSheet,
} from 'react-native';
import TeaListItem from './TeaListItem';
import AddTea from './AddTea';
import {serverURL} from '../../app.json';

const TeaManager = () => {
  const [teasByType, setTeasByType] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  const deleteTea = id => {
    fetch(serverURL + '/tea/' + id, {
      method: 'DELETE',
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        } else {
          let newTeasByType = [];
          teasByType.forEach(teaByType => {
            let newTeaByType = {title: '', data: []};
            newTeaByType.title = teaByType.title;
            teaByType.data.forEach(tea => {
              if (tea.id !== id) {
                newTeaByType.data.push(tea);
              }
            });
            newTeasByType.push(newTeaByType);
          });
          setTeasByType(newTeasByType);
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
    } else if (
      teasByType.some(teas => teas.data.some(tea => tea.name === name))
    ) {
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
        let newTeasByType = [...teasByType];
        let index = newTeasByType.findIndex(
          teaType => teaType.title === 'Black Tea',
        ); ///TODO: Add in type selected by user
        if (index === -1) {
          newTeasByType.push({title: 'Black Tea', data: [{id: json.id, name}]});
        } else {
          newTeasByType[index].data.push({id: json.id, name});
        }
        setTeasByType(newTeasByType);
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
    fetch(serverURL + '/types/teas')
      .then(response => response.json())
      .then(json => {
        let responseTeasByType = [];
        json.forEach(teaByType => {
          let newTeasByType = {title: '', data: []};
          newTeasByType.title = teaByType.type.name;
          teaByType.teas.forEach(tea => {
            newTeasByType.data.push({id: tea.id, name: tea.name});
          });
          responseTeasByType.push(newTeasByType);
        });
        setTeasByType(responseTeasByType);
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
          sections={teasByType}
          renderItem={({item}) => (
            <TeaListItem item={item} deleteTea={deleteTea} />
          )}
          renderSectionHeader={({section: {title}}) => (
            <Text style={styles.header}>{title}</Text>
          )}
          keyExtractor={(item, index) => item + index}
          ListFooterComponent={<AddTea addTea={addTea} />}
          ref={ref => (this.sectionListRef = ref)}
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  header: {
    fontSize: 24,
    backgroundColor: '#fff',
    paddingTop: 20,
    paddingBottom: 5,
    textAlign: 'center',
    fontWeight: '900',
  },
});

export default TeaManager;
