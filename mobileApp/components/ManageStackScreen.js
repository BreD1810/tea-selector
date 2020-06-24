import React from 'react';
import {Text} from 'react-native';
import {createStackNavigator} from '@react-navigation/stack';
import ManageButtons from './ManageButtons';
import TeaManager from './TeaManager/TeaManager';
import TeaTypeManager from './TeaManager/TeaTypeManager';

const ManageStack = createStackNavigator();

const TempComponent = () => {
  return <Text>Nothing here yet!</Text>;
};

const ManageStackScreen = () => {
  return (
    <ManageStack.Navigator initialRouteName="Manage">
      <ManageStack.Screen
        options={{
          headerShown: false,
        }}
        name="Manage"
        component={ManageButtons}
      />
      <ManageStack.Screen
        name="TeaManager"
        component={TeaManager}
        options={{
          headerTitle: 'Manage Teas',
        }}
      />
      <ManageStack.Screen
        name="TeaTypeManager"
        component={TeaTypeManager}
        options={{
          headerTitle: 'Manage Tea Types',
        }}
      />
      <ManageStack.Screen name="Temp" component={TempComponent} />
    </ManageStack.Navigator>
  );
};

export default ManageStackScreen;
