import React from 'react';
import {createStackNavigator} from '@react-navigation/stack';
import ManageButtons from './ManageButtons';
import TeaManager from './Managers/TeaManager';
import TeaTypeManager from './Managers/TeaTypeManager';
import OwnerManager from './Managers/OwnerManager';
import OwnershipManager from './Managers/OwnershipManager';
import AccountManager from './Managers/AccountManager';

const ManageStack = createStackNavigator();

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
      <ManageStack.Screen
        name="AccountManager"
        children={() => <AccountManager jwtToken={jwtToken} />}
        options={{
          headerTitle: 'Manage Account',
        }}
      />
    </ManageStack.Navigator>
  );
};

export default ManageStackScreen;
