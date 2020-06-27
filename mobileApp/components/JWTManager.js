import {AsyncStorage} from 'react-native';
import {serverURL} from '../app.json';

const JWTManager = {
  async getJWT(setJWT, setIsAuthorized, setIsLoading) {
    try {
      const value = await AsyncStorage.getItem('jwt_token');
      if (value !== null) {
        console.log('Found JWT');
        fetch(serverURL + '/owners', {
          headers: {
            Token: value,
          },
        })
          .then(response => {
            if (!response.ok) {
              console.log('No longer valid');
              throw new Error(response.statusText);
            }
            console.log('Still valid');
            setJWT(value);
            setIsAuthorized(true);
          })
          .catch(error => {
            console.warn(error);
          })
          .finally(() => setIsLoading(false));
      } else {
        console.log('No JWT');
        setIsLoading(false);
      }
    } catch (error) {
      console.warn('AsyncStorage error: ' + error.message);
      setIsLoading(false);
    }
  },
  async setJWT(jwt) {
    try {
      await AsyncStorage.setItem('jwt_token', jwt);
    } catch (error) {
      console.warn('AsyncStorage error: ' + error.message);
    }
  },
  async clearJWT() {
    try {
      await AsyncStorage.removeItem('jwt_token');
    } catch (error) {
      console.warn('AsyncStorage error: ' + error.message);
    }
  },
};

export default JWTManager;
