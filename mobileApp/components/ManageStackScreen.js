import React from 'react';
import {Text} from 'react-native';
import {createStackNavigator} from '@react-navigation/stack';
import ManageButtons from './ManageButtons';
import TeaManager from './Managers/TeaManager';
import TeaTypeManager from './Managers/TeaTypeManager';
import OwnerManager from './Managers/OwnerManager';
import OwnershipManager from './Managers/OwnershipManager';

const ManageStack = createStackNavigator();

const TempComponent = () => {
  return <Text>Nothing here yet!</Text>;
};

const ManageStackScreen = ({jwtToken}) => {
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
        children={() => <TeaManager jwtToken={jwtToken} />}
        options={{
          headerTitle: 'Manage Teas',
        }}
      />
      <ManageStack.Screen
        name="TeaTypeManager"
        children={() => <TeaTypeManager jwtToken={jwtToken} />}
        options={{
          headerTitle: 'Manage Tea Types',
        }}
      />
      <ManageStack.Screen
        name="OwnerManager"
        children={() => <OwnerManager jwtToken={jwtToken} />}
        options={{
          headerTitle: 'Manage Owners',
        }}
      />
      <ManageStack.Screen
        name="OwnershipManager"
        children={() => <OwnershipManager jwtToken={jwtToken} />}
        options={{
          headerTitle: 'Manage Tea Owners',
        }}
      />
    </ManageStack.Navigator>
  );
};

export default ManageStackScreen;
