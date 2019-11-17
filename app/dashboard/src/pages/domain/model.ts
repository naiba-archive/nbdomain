import { AnyAction, Reducer } from 'redux';
import { EffectsCommandMap } from 'dva';
// eslint-disable-next-line
import { TableListData } from './data';

import { APIS } from '@/services';

export interface StateType {
  data?: TableListData;
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
    whois: Effect;
  };
  reducers: {
    save: Reducer<StateType>;
  };
}

const Model: ModelType = {
  namespace: 'domain',

  state: {
    data: {
      list: [],
      pagination: {},
    },
  },

  effects: {
    *fetch({ payload }, { call, put }) {
      const response = yield call(APIS.DefaultApi.domainGet, payload);
      yield put({
        type: 'save',
        payload: response,
      });
    },
    *add({ payload, callback }, { call, put }) {
      const response = yield call(APIS.DefaultApi.domainPost, payload);
      yield put({
        type: 'save',
        payload: response,
      });
      if (callback && response) callback();
    },
    *remove({ payload, callback }, { call }) {
      const response = yield call(APIS.DefaultApi.domainIdDelete, payload);
      if (callback && response) callback();
    },
    *whois({ payload, callback }, { call }) {
      const response = yield call(APIS.DefaultApi.whoisDomainGet, payload);
      if (callback && response) callback(response);
    },
  },

  reducers: {
    save(state, action) {
      return {
        ...state,
        data: action.payload,
      };
    },
  },
};

export default Model;
