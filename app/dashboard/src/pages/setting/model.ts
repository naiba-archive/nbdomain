import { AnyAction, Reducer } from 'redux';
import { EffectsCommandMap } from 'dva';
import { APIS } from '@/services';

export type Effect = (
  action: AnyAction,
  effects: EffectsCommandMap & { select: <T>(func: (state: {}) => T) => T },
) => void;

export interface SettingState {
  currentUser: any;
}

export interface ModelType {
  namespace: string;
  state: SettingState;
  effects: {
    fetchCurrent: Effect;
    submitRegularForm: Effect;
  };
  reducers: {
    save: Reducer<SettingState>;
  };
}

const Model: ModelType = {
  namespace: 'setting',

  state: {
    currentUser: {},
  },

  effects: {
    *submitRegularForm({ payload, callback }, { call }) {
      const response = yield call(APIS.DefaultApi.userPut, payload);
      if (response && callback) callback(response);
    },
    *fetchCurrent(_, { call, put }) {
      const response = yield call(APIS.DefaultApi.userGet);
      yield put({
        type: 'save',
        payload: response,
      });
    },
  },

  reducers: {
    save(state, action) {
      return {
        ...state,
        currentUser: action.payload,
      };
    },
  },
};

export default Model;
