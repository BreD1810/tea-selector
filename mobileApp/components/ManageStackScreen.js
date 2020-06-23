import React from 'react';
import {
  Text,
} from 'react-native';
import { createStackNavigator } from '@react-navigation/stack';
import ManageButtons from './ManageButtons'

const ManageStack = createStackNavigator();

const Test2Component = () => {
  return <Text>Manage</Text>
}

const ManageStackScreen = (props) => {
  return (
    <ManageStack.Navigator initialRouteName="Manage">
      <ManageStack.Screen
        options={{
          headerShown: false
        }}
        name="Manage"
        component={ManageButtons} 
      />
      <ManageStack.Screen name="Test2" component={Test2Component} />
    </ManageStack.Navigator>
  );
};

export default ManageStackScreen;
