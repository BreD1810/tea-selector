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

const OwnershipManager = ({jwtToken}) => {
  const [ownersTeas, setOwnersTeas] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [teas, setTeas] = useState([]);

  const deleteTea = (teaID, ownerID) => {
    fetch(serverURL + '/tea/' + teaID + '/owner/' + ownerID, {
      method: 'DELETE',
      headers: {
        Token: jwtToken,
      },
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
        Token: jwtToken,
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
        newOwnersTeas[index].data.sort((a, b) =>
          a.name > b.name ? 1 : b.name > a.name ? -1 : 0,
        ); // Sort teas alphabetically
        setOwnersTeas(newOwnersTeas);
        ToastAndroid.showWithGravityAndOffset(
          'Tea successfully added!',
          ToastAndroid.SHORT,
          ToastAndroid.BOTTOM,
          0,
          150,
        );
        resetFunc();
      })
      .catch(error => {
        console.warn(error);
        Alert.alert('Error adding tea', 'Please try again.');
      });
  };

  const getTeas = () => {
    let newTeas = [];
    fetch(serverURL + '/teas', {
      headers: {
        Token: jwtToken,
      },
    })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        }
        return response.json();
      })
      .then(json => {
        json.forEach(tea => newTeas.push({label: tea.name, value: tea.id}));
      })
      .then(() => {
        newTeas = newTeas.sort((a, b) =>
          a.label > b.label ? 1 : b.label > a.label ? -1 : 0,
        ); // Sort alphabetically
      })
      .then(() => setTeas(newTeas))
      .catch(error => {
        console.warn(error);
      });
  };

  useEffect(() => {
    fetch(serverURL + '/owners/teas', {
      headers: {
        Token: jwtToken,
      },
    })
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
            newOwnerTeas.data.sort((a, b) =>
              a.name > b.name ? 1 : b.name > a.name ? -1 : 0,
            ); // Sort each person's teas alphabetically.
          }
          newOwnersTeas.push(newOwnerTeas);
        });
        newOwnersTeas.sort((a, b) =>
          a.title.name > b.title.name
            ? 1
            : b.title.name > a.title.name
            ? -1
            : 0,
        ); // Sort owners alphabetically.
        setOwnersTeas(newOwnersTeas);
      })
      .then(() => getTeas())
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
            <AddSectionItemPicker
              promptText={'Select A Tea...'}
              inputItems={teas}
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
