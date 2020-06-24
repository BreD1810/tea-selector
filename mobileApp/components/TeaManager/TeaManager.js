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
import ListItem from './ListItem';
import AddSectionItem from './AddSectionItem';
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

  const addTea = (name, typeID, textInput) => {
    if (!name) {
      Alert.alert('Error', 'Please enter a name for the new tea!');
      return;
    } else if (
      teasByType.some(teas => teas.data.some(tea => tea.name === name))
    ) {
      Alert.alert('Error', 'That tea already exists!');
      return;
    }

    addTeaAPI(name, typeID, textInput);
  };

  const addTeaAPI = (name, typeID, textInput) => {
    fetch(serverURL + '/tea', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({name: name, type: {id: typeID}}),
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
          teaType => teaType.title.id === typeID,
        );
        newTeasByType[index].data.push({id: json.id, name});
        setTeasByType(newTeasByType);
        ToastAndroid.showWithGravityAndOffset(
          'Tea successfully added!',
          ToastAndroid.SHORT,
          ToastAndroid.BOTTOM,
          0,
          150,
        );
        textInput.clear();
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
          newTeasByType.title = {
            id: teaByType.type.id,
            name: teaByType.type.name,
          };
          if (teaByType.teas !== null) {
            teaByType.teas.forEach(tea => {
              newTeasByType.data.push({id: tea.id, name: tea.name});
            });
          }
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
          style={styles.sectionList}
          sections={teasByType}
          renderItem={({item}) => (
            <ListItem item={item} deleteFunc={deleteTea} />
          )}
          renderSectionHeader={({section: {title}}) => (
            <Text style={styles.header}>{title.name}</Text>
          )}
          keyExtractor={(item, index) => item + index}
          renderSectionFooter={({section: {title}}) => (
            <AddSectionItem
              placeholderText={'Add Tea...'}
              addFunc={addTea}
              sectionID={title.id}
            />
          )}
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  header: {
    fontSize: 24,
    paddingTop: 10,
    paddingBottom: 5,
    textAlign: 'center',
    fontWeight: '900',
  },
});

export default TeaManager;
