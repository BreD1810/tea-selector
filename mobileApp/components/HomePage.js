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
    } else {
      setSelectedTea(Math.floor(Math.random() * teas.length));
    }
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
    if (checkboxStates.some(state => state)) {
      getAllTeasRespectingOwners();
    } else {
      getAllTeas();
    }
  };

  const getAllTeas = () => {
    fetch(serverURL + '/teas')
      .then(response => response.json())
      .then(json => {
        setIsLoading(true);
        setTeas(json);
        setSelectedTea(Math.floor(Math.random() * json.length));
      })
      .catch(error => console.error(error))
      .finally(() => {
        setIsLoading(false);
      });
  };

  const getAllTeasRespectingOwners = () => {
    fetch(serverURL + '/owners/teas')
      .then(response => {
        if (!response.ok) {
          throw new Error(response.json().error);
        }
        return response.json();
      })
      .then(json => {
        setIsLoading(true);
        let ownerIDs = [];
        let newTeas = [];
        checkboxStates.forEach((state, index) => {
          if (state) {
            ownerIDs.push(owners[index].id);
          }
        });
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
        newTeas = newTeas.filter(tea => tea.count === ownerIDs.length);
        setTeas(newTeas);
        setSelectedTea(Math.floor(Math.random() * newTeas.length));
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
        ) : teas[selectedTea] === undefined ? (
          <Text style={styles.tea}>Error... Please try again</Text>
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
      <Text style={styles.subtitle}>Filter by owners:</Text>
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
    fontSize: 22,
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
