import React from 'react';
import 'react-native-gesture-handler';
import { Text } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import Icon from 'react-native-vector-icons/dist/FontAwesome5'
import HomePage from './components/HomePage';

const ManageComponent = () => {
  return <Text>Manage</Text>
}

const Tab = createBottomTabNavigator();

const App: () => React$Node = () => {
  return (
    <NavigationContainer>
      <Tab.Navigator
        screenOptions={({ route }) => ({
          tabBarIcon: ({ focused, color, size }) => {
            let iconName;

            if (route.name === 'Home') {
              iconName = focused ? 'home' : 'home';
            } else if (route.name === 'Manage') {
              iconName = focused ? 'sliders-h' : 'sliders-h';
            }

            return <Icon name={iconName} size={size} color={color} />;
          },
        })}
        tabBarOptions = {{
          activeTintColor: 'dodgerblue',
          inactiveTintColor: 'grey',
        }}
      >
      <Tab.Screen name="Home" component={HomePage} />
      <Tab.Screen name="Manage" component={ManageComponent} />
      </Tab.Navigator>
    </NavigationContainer>
  );
};


export default App;
