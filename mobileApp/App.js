import React, {useState, useEffect} from 'react';
import 'react-native-gesture-handler';
import {NavigationContainer} from '@react-navigation/native';
import {createBottomTabNavigator} from '@react-navigation/bottom-tabs';
import Icon from 'react-native-vector-icons/dist/FontAwesome5';
import HomePage from './components/HomePage';
import ManageStackScreen from './components/ManageStackScreen';
import JWTManager from './components/JWTManager';
import Login from './components/Login';

const Tab = createBottomTabNavigator();

const App: () => React$Node = () => {
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [jwt, setJWT] = useState(null);

  const login = newJWT => {
    setJWT(newJWT);
    JWTManager.setJWT(newJWT);
    setIsAuthorized(true);
  };

  useEffect(() => {
    checkAuthorized();
  }, []);

  const checkAuthorized = () => {
    JWTManager.getJWT(setJWT, setIsAuthorized);
  };

  return (
    <>
      {isAuthorized ? (
        <NavigationContainer>
          <Tab.Navigator
            screenOptions={({route}) => ({
              tabBarIcon: ({focused, color, size}) => {
                let iconName;

                if (route.name === 'Home') {
                  iconName = focused ? 'home' : 'home';
                } else if (route.name === 'Manage') {
                  iconName = focused ? 'sliders-h' : 'sliders-h';
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
              name="Manage"
              children={() => <ManageStackScreen jwtToken={jwt} />}
            />
          </Tab.Navigator>
        </NavigationContainer>
      ) : (
        <Login setJWT={login} />
      )}
    </>
  );
};

export default App;
