import React, {useState, useEffect} from 'react';
import {serverURL} from '../app.json';
import {View, Text, StyleSheet, ActivityIndicator} from 'react-native';
import Button from 'react-native-button';
import CheckBoxGroup from './CheckBoxGroup';

const HomePage = () => {
  const [teas, setTeas] = useState([]);
  const [selectedTea, setSelectedTea] = useState(0);
  const [owners, setOwners] = useState([]);
  const [checkboxStates, setCheckboxStates] = useState([]);
  const [checkboxesChanged, setCheckboxesChanged] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingOwners, setIsLoadingOwners] = useState(true);

  const selectRandomTea = () => {
    if (checkboxesChanged) {
      setCheckboxesChanged(false);
      updateTeas();
    }
    setSelectedTea(Math.floor(Math.random() * teas.length));
  };

  const checkboxChange = index => {
    let newCheckboxStates = [...checkboxStates];
    newCheckboxStates[index] = !checkboxStates[index];
    setCheckboxStates(newCheckboxStates);
    setCheckboxesChanged(true);
  };

  useEffect(() => {
    getTeas();
    getOwners();
  }, []);

  const getTeas = () => {
    console.log('Getting teas');
    fetch(serverURL + '/teas')
      .then(response => response.json())
      .then(json => {
        setTeas(json);
      })
      .then(() => selectRandomTea())
      .catch(error => console.error(error))
      .finally(() => {
        setIsLoading(false);
      });
  };

  const updateTeas = () => {
    console.log('Updating teas');
    if (checkboxStates.some(state => state)) {
      console.log('Some box checked');
      getAllTeasRespectingOwners();
    } else {
      console.log('Nothing checked');
      getAllTeas();
    }
  };

  const getAllTeas = () => {
    setIsLoading(true);
    fetch(serverURL + '/teas')
      .then(response => response.json())
      .then(json => {
        setTeas(json);
      })
      .catch(error => console.error(error))
      .finally(() => {
        setIsLoading(false);
      });
  };

  const getAllTeasRespectingOwners = () => {
    setIsLoading(true);
    let ownerIDs = [];
    checkboxStates.forEach((state, index) => {
      if (state) {
        ownerIDs.push(owners[index].id);
      }
    });
    fetch(serverURL + '/owners/teas')
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        }
        return response.json();
      })
      .then(json => {
        let newTeas = [];
        console.log(json);
        json.forEach(ownerTeasResponse => {
          if (ownerIDs.includes(ownerTeasResponse.owner.id)) {
            let ownersTeas = ownerTeasResponse.teas;
            ownersTeas.forEach(tea => {
              let index = newTeas.findIndex(otherTea => otherTea.id === tea.id);
              if (index === -1) {
                newTeas.push({...tea, count: 1});
              } else {
                newTeas[index].count = newTeas[index].count + 1;
              }
            });
          }
        });
        console.log(ownerIDs.length);
        console.log(newTeas);
        newTeas = newTeas.filter(tea => tea.count === ownerIDs.length);
        console.log('Before: ' + teas);
        setTeas(newTeas);
        console.log('After: ' + teas);
      })
      .catch(error => console.error(error))
      .finally(() => {
        setIsLoading(false);
      });
  };

  const getOwners = () => {
    fetch(serverURL + '/owners')
      .then(response => response.json())
      .then(json => {
        setOwners(json);
      })
      .catch(error => console.error(error))
      .finally(() => {
        setIsLoadingOwners(false);
      });
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>The selected tea is:</Text>
      <View style={styles.teaTextContainer}>
        {isLoading || selectedTea === null ? (
          <ActivityIndicator />
        ) : teas === [] ? (
          <Text style={styles.tea}>No teas in common!</Text>
        ) : (
          <Text style={styles.tea}>{teas[selectedTea].name}</Text>
        )}
      </View>
      <Button
        onPress={selectRandomTea}
        containerStyle={styles.btnContainer}
        style={styles.btn}>
        Select Another Tea
      </Button>
      <Text style={styles.subtitle}>Owners</Text>
      {isLoadingOwners ? (
        <ActivityIndicator />
      ) : (
        <CheckBoxGroup
          items={owners}
          states={checkboxStates}
          updateFunc={checkboxChange}
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
  },
  teaTextContainer: {
    height: 40,
  },
  title: {
    fontSize: 32,
    fontWeight: '600',
    color: 'black',
  },
  subtitle: {
    paddingTop: 25,
    fontSize: 30,
    fontWeight: '600',
    color: 'black',
  },
  tea: {
    marginTop: 8,
    fontSize: 24,
    fontWeight: '400',
    color: 'black',
  },
  btn: {
    fontSize: 20,
    color: 'white',
  },
  btnContainer: {
    padding: 15,
    marginTop: 20,
    height: 60,
    width: 200,
    overflow: 'hidden',
    borderRadius: 4,
    backgroundColor: 'dodgerblue',
  },
});

export default HomePage;
