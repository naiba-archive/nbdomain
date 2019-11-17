import { Button, Card, Form, Input } from 'antd';
import React, { Component } from 'react';

import { Dispatch } from 'redux';
import { FormComponentProps } from 'antd/es/form';
import { PageHeaderWrapper } from '@ant-design/pro-layout';
import { connect } from 'dva';
import { SettingState } from './model';

const FormItem = Form.Item;

interface SettingProps extends FormComponentProps {
  setting: SettingState;
  submitting: boolean;
  dispatch: Dispatch<any>;
}

class Setting extends Component<SettingProps> {
  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({
      type: 'setting/fetchCurrent',
    });
  }

  handleSubmit = (e: React.FormEvent) => {
    const { dispatch, form } = this.props;
    e.preventDefault();
    form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        dispatch({
          type: 'setting/submitRegularForm',
          payload: values,
          callback: () => {
            window.location.reload();
          },
        });
      }
    });
  };

  render() {
    const {
      submitting,
      setting: { currentUser },
    } = this.props;
    const {
      form: { getFieldDecorator },
    } = this.props;

    const formItemLayout = {
      labelCol: {
        xs: { span: 24 },
        sm: { span: 7 },
      },
      wrapperCol: {
        xs: { span: 24 },
        sm: { span: 12 },
        md: { span: 10 },
      },
    };

    const submitFormLayout = {
      wrapperCol: {
        xs: { span: 24, offset: 0 },
        sm: { span: 10, offset: 7 },
      },
    };

    return (
      <PageHeaderWrapper>
        <Card bordered={false}>
          <Form onSubmit={this.handleSubmit} hideRequiredMark style={{ marginTop: 8 }}>
            <FormItem {...formItemLayout} label="昵称">
              {getFieldDecorator('name', {
                initialValue: currentUser.name,
              })(<Input placeholder="Leo" />)}
            </FormItem>
            <FormItem {...formItemLayout} label="手机号">
              {getFieldDecorator('phone', {
                initialValue: currentUser.phone,
              })(<Input placeholder="18888888888" />)}
            </FormItem>
            <FormItem {...formItemLayout} label="微信">
              {getFieldDecorator('weixin', {
                initialValue: currentUser.weixin,
              })(<Input placeholder="Leo88888" />)}
            </FormItem>
            <FormItem {...formItemLayout} label="QQ">
              {getFieldDecorator('qq', {
                initialValue: currentUser.qq,
              })(<Input placeholder="88888" />)}
            </FormItem>
            <FormItem {...formItemLayout} label="新密码">
              {getFieldDecorator('password', {})(
                <Input type="password" placeholder="留空不修改" />,
              )}
            </FormItem>
            <FormItem {...submitFormLayout} style={{ marginTop: 32 }}>
              <Button type="primary" htmlType="submit" loading={submitting}>
                保存
              </Button>
            </FormItem>
          </Form>
        </Card>
      </PageHeaderWrapper>
    );
  }
}

export default Form.create<SettingProps>()(
  connect(
    ({
      setting,
      loading,
    }: {
      setting: SettingState;
      loading: { effects: { [key: string]: boolean } };
    }) => ({
      setting,
      submitting: loading.effects['setting/submitRegularForm'],
    }),
  )(Setting),
);
