import React, {useState, useEffect} from 'react';
import 'react-native-gesture-handler';
import {NavigationContainer} from '@react-navigation/native';
import {createBottomTabNavigator} from '@react-navigation/bottom-tabs';
import Icon from 'react-native-vector-icons/dist/FontAwesome5';
import HomePage from './components/HomePage';
import ManageStackScreen from './components/ManageStackScreen';
import JWTManager from './components/JWTManager';
import Login from './components/Login';
import AccountManager from './components/Managers/AccountManager';
import {ActivityIndicator, StyleSheet, View} from 'react-native';

const Tab = createBottomTabNavigator();

const App: () => React$Node = () => {
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [jwt, setJWT] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const login = newJWT => {
    setJWT(newJWT);
    JWTManager.setJWT(newJWT);
    setIsAuthorized(true);
  };

  useEffect(() => {
    JWTManager.getJWT(setJWT, setIsAuthorized, setIsLoading);
  }, []);

  return (
    <>
      {isLoading ? (
        <View style={styles.container}>
          <ActivityIndicator size={100} color="dodgerblue" />
        </View>
      ) : isAuthorized ? (
        <NavigationContainer>
          <Tab.Navigator
            screenOptions={({route}) => ({
              tabBarIcon: ({focused, color, size}) => {
                let iconName;

                if (route.name === 'Home') {
                  iconName = 'home';
                } else if (route.name === 'Manage Tea') {
                  iconName = 'coffee';
                } else if (route.name === 'Account') {
                  iconName = 'user';
                }

                return <Icon name={iconName} size={size} color={color} />;
              },
            })}
            tabBarOptions={{
              activeTintColor: 'dodgerblue',
              inactiveTintColor: 'grey',
            }}>
            <Tab.Screen
              name="Home"
              children={() => <HomePage jwtToken={jwt} />}
            />
            <Tab.Screen
              name="Manage Tea"
              children={() => <ManageStackScreen jwtToken={jwt} />}
            />
            <Tab.Screen
              name="Account"
              children={() => <AccountManager jwtToken={jwt} setLoggedIn={setIsAuthorized}/>}
            />
          </Tab.Navigator>
        </NavigationContainer>
      ) : (
        <Login setJWT={login} />
      )}
    </>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
  },
});
export default App;
