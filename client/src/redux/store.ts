import {combineReducers, createStore, applyMiddleware, compose} from "redux";
import reduxThunk from "redux-thunk"
import {credentialsReducer} from "./credential";

declare var window: {
    __REDUX_DEVTOOLS_EXTENSION_COMPOSE__: any,
};


const reducers = combineReducers({
    CredentialsReducer: credentialsReducer
});

// Add redux-dev-tools on client side only middleware
function CreateStore(preloadState?: any) {
    if (typeof window !== "undefined") {
        console.log('Client side store setup');
        const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
        return createStore(reducers, preloadState, composeEnhancers(applyMiddleware(reduxThunk)));
    }

    console.log('Server side store setup');
    return createStore(reducers, preloadState, applyMiddleware(reduxThunk));
}

export default CreateStore;