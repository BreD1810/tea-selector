Tea Selector
---

I can never decide which tea to have, and neither can my friends.
This app choose the tea for us!

## Mobile App

The mobile app was built using [React Native](https://reactnative.dev/), targeting Android.

The main function of the app is to select a tea that each 'drinker' posses.
If none are selected, a random tea is selected.

<img src="imgs/main.jpg" height="500">

Beyond that, the user is free to use a range of management options.
This includes adding/deleting teas, adding/deleting tea types, adding/deleting owners, and also adding/deleting teas from users.

<img src="imgs/management-tab.jpg" height="500">
<img src="imgs/tea-manager.jpg" height="500">
<img src="imgs/ownership.jpg" height="500">

Additionally, a login is required to use the app.
The 'account' tab lets the user change their password, or log out.

<img src="imgs/accounts.jpg" height="500">

## API
[![pipeline status](https://gitlab.com/BreD1810/tea-selector/badges/master/pipeline.svg)](https://gitlab.com/BreD1810/tea-selector/-/commits/master)
[![coverage report](https://gitlab.com/BreD1810/tea-selector/badges/master/coverage.svg)](https://gitlab.com/BreD1810/tea-selector/-/commits/master)

This API provides all the data and functionality needed for the app.
It was built using [Golang](https://golang.org/).
You can read more about the API [here](api/README.md)
