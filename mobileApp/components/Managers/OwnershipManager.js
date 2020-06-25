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
import ListItem from './Lists/ListItem';
import {serverURL} from '../../app.json';
import AddSectionItemPicker from './Lists/AddSectionItemPicker';

const OwnershipManager = () => {
  const [ownersTeas, setOwnersTeas] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  const deleteTea = (teaID, ownerID) => {
    fetch(serverURL + '/tea/' + teaID + '/owner/' + ownerID, {
      method: 'DELETE',
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        } else {
          let newOwnersTeas = [];
          ownersTeas.forEach(teaOwner => {
            if (teaOwner.title.id !== ownerID) {
              newOwnersTeas.push(teaOwner);
            } else {
              let newOwnerTeas = {title: '', data: []};
              newOwnerTeas.title = {
                id: teaOwner.title.id,
                name: teaOwner.title.name,
              };
              teaOwner.data.forEach(tea => {
                if (tea.id !== teaID) {
                  newOwnerTeas.data.push(tea);
                }
              });
              newOwnersTeas.push(newOwnerTeas);
            }
          });
          setOwnersTeas(newOwnersTeas);
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
        Alert.alert('Error deleting tea', 'Not sure what went wrong here...');
      });
  };

  const addOwnersTea = (teaID, ownerID, resetFunc) => {
    if (teaID === null) {
      Alert.alert('Error', 'Please select a tea!');
      return;
    } else {
      let index = ownersTeas.findIndex(owner => owner.title.id === ownerID);
      if (ownersTeas[index].data.some(tea => tea.id === teaID)) {
        Alert.alert('Error', 'That person already owns that tea!');
        return;
      }
    }

    addOwnersTeaAPI(teaID, ownerID, resetFunc);
  };

  const addOwnersTeaAPI = (teaID, ownerID, resetFunc) => {
    fetch(serverURL + '/tea/' + teaID + '/owner', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({id: ownerID}),
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        }
        return response.json();
      })
      .then(json => {
        let newOwnersTeas = [...ownersTeas];
        let index = newOwnersTeas.findIndex(
          owner => owner.title.id === ownerID,
        );
        newOwnersTeas[index].data.push({id: json.id, name: json.name});
        setOwnersTeas(newOwnersTeas);
        ToastAndroid.showWithGravityAndOffset(
          'Tea successfully added!',
          ToastAndroid.SHORT,
          ToastAndroid.BOTTOM,
          0,
          150,
        );
        resetFunc();
        // textInput.clear();
      })
      .catch(error => {
        console.warn(error);
        Alert.alert('Error adding tea', 'Please try again.');
      });
  };

  const getTeas = () => {
    let teas = [];
    fetch(serverURL + '/teas')
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        }
        return response.json();
      })
      .then(json => {
        json.forEach(tea => teas.push({label: tea.name, value: tea.id}));
      })
      .catch(error => {
        console.warn(error);
      });
    return teas;
  };

  useEffect(() => {
    fetch(serverURL + '/owners/teas')
      .then(response => response.json())
      .then(json => {
        let newOwnersTeas = [];
        json.forEach(teaOwner => {
          let newOwnerTeas = {title: '', data: []};
          newOwnerTeas.title = {
            id: teaOwner.owner.id,
            name: teaOwner.owner.name,
          };
          if (teaOwner.teas !== null) {
            teaOwner.teas.forEach(tea => {
              newOwnerTeas.data.push({id: tea.id, name: tea.name});
            });
          }
          newOwnersTeas.push(newOwnerTeas);
        });
        setOwnersTeas(newOwnersTeas);
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
          sections={ownersTeas}
          renderItem={({item, section}) => (
            <ListItem
              item={item}
              deleteFunc={deleteTea}
              sectionID={section.title.id}
            />
          )}
          renderSectionHeader={({section: {title}}) => (
            <Text style={styles.header}>{title.name}</Text>
          )}
          keyExtractor={(item, index) => item + index}
          renderSectionFooter={({section: {title}}) => (
            <AddSectionItemPicker // TODO: Change this to be a selector for teas?
              promptText={'Select A Tea...'}
              inputItems={getTeas()}
              addFunc={addOwnersTea}
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

export default OwnershipManager;
