import { AnyAction, Reducer } from 'redux';
import { EffectsCommandMap } from 'dva';
import { TableListData } from './data.d';

import { APIS } from '@/services';

export interface StateType {
  data?: TableListData;
  panelOptions?: any;
}

export type Effect = (
  action: AnyAction,
  effects: EffectsCommandMap & { select: <T>(func: (state: StateType) => T) => T },
) => void;

export interface ModelType {
  namespace: string;
  state: StateType;
  effects: {
    fetch: Effect;
    fetchOptions: Effect;
    add: Effect;
    remove: Effect;
  };
  reducers: {
    save: Reducer<StateType>;
    saveOptions: Reducer<StateType>;
  };
}

const Model: ModelType = {
  namespace: 'panel',

  state: {
    data: {
      list: [],
      pagination: {},
    },
    panelOptions: {},
  },

  effects: {
    *fetch({ payload }, { call, put }) {
      const response = yield call(APIS.DefaultApi.panelGet, payload);
      yield put({
        type: 'save',
        payload: response,
      });
    },
    *fetchOptions({ payload }, { call, put }) {
      const response = yield call(APIS.DefaultApi.panelOptionGet, payload);
      yield put({
        type: 'saveOptions',
        payload: response,
      });
    },
    *add({ payload, callback }, { call, put }) {
      const response = yield call(APIS.DefaultApi.panelPost, payload);
      yield put({
        type: 'save',
        payload: response,
      });
      if (callback && response) callback();
    },
    *remove({ payload, callback }, { call }) {
      const response = yield call(APIS.DefaultApi.panelidDelete, payload);
      if (callback && response) callback();
    },
  },

  reducers: {
    save(state, action) {
      return {
        ...state,
        data: action.payload,
      };
    },
    saveOptions(state, action) {
      return {
        ...state,
        panelOptions: action.payload,
      };
    },
  },
};

export default Model;
