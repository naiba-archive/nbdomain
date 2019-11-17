import { AnyAction, Reducer } from 'redux';
import { EffectsCommandMap } from 'dva';
import { TableListData } from './data.d';

import { APIS } from '@/services';

export interface StateType {
  data?: TableListData;
  catOptions?: any;
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
    add: Effect;
    remove: Effect;
  };
  reducers: {
    save: Reducer<StateType>;
    saveOptions: Reducer<StateType>;
  };
}

const Model: ModelType = {
  namespace: 'cat',

  state: {
    data: {
      list: [],
      pagination: {},
    },
    catOptions: {},
  },

  effects: {
    *fetch({ payload }, { call, put }) {
      const response = yield call(APIS.DefaultApi.catGet, payload);
      yield put({
        type: 'save',
        payload: response,
      });
    },
    *add({ payload, callback }, { call, put }) {
      const response = yield call(APIS.DefaultApi.catPost, payload);
      yield put({
        type: 'save',
        payload: response,
      });
      if (callback && response) callback();
    },
    *remove({ payload, callback }, { call }) {
      const response = yield call(APIS.DefaultApi.catidDelete, payload);
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
        catOptions: action.payload,
      };
    },
  },
};

export default Model;
