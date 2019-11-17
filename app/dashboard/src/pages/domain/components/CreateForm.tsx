import { Form, Input, Modal, InputNumber, DatePicker } from 'antd';

import { FormComponentProps } from 'antd/es/form';
import React from 'react';
import TextArea from 'antd/lib/input/TextArea';
import moment from 'moment';

const FormItem = Form.Item;

interface CreateFormProps extends FormComponentProps {
  dispatch: any;
  isEdit: boolean;
  modalVisible: boolean;
  currentRow: any;
  handleAdd: (fieldsValue: any, isEdit: boolean) => void;
  handleModalVisible: () => void;
}

const CreateForm: React.FC<CreateFormProps> = props => {
  const { modalVisible, isEdit, dispatch, currentRow, form, handleAdd, handleModalVisible } = props;

  const okHandle = () => {
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      handleAdd(fieldsValue, isEdit);
    });
  };

  const handleBlur = () => {
    dispatch({
      type: 'domain/whois',
      payload: { domain: form.getFieldValue('domain') },
      callback: (resp: any) => {
        form.setFieldsValue({
          registrar: resp.registrar.name,
          create: moment(resp.domain.created_date),
          expire: moment(resp.domain.expiration_date),
        });
      },
    });
  };

  return (
    <Modal
      destroyOnClose
      title={`${isEdit ? '修改' : '新建'}域名`}
      visible={modalVisible}
      onOk={okHandle}
      onCancel={() => handleModalVisible()}
    >
      {form.getFieldDecorator('id', {
        initialValue: currentRow.id ? currentRow.id : 0,
      })(<input type="hidden" />)}
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="米表ID">
        {form.getFieldDecorator('panel_id', {
          rules: [{ required: true, message: '请输入米表ID' }],
          initialValue: currentRow.panel_id ? currentRow.panel_id : 1,
        })(<InputNumber min={1} placeholder="1" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="分类ID">
        {form.getFieldDecorator('cat_id', {
          rules: [{ required: true, message: '请输入米表ID' }],
          initialValue: currentRow.cat_id ? currentRow.cat_id : 1,
        })(<InputNumber min={1} placeholder="1" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="域名">
        {form.getFieldDecorator('domain', {
          rules: [{ required: true, message: '请输入域名' }],
          initialValue: currentRow.domain ? currentRow.domain : '',
        })(<Input onBlur={handleBlur} placeholder="example.com" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="简介">
        {form.getFieldDecorator('desc', {
          rules: [{ required: true, message: '请输入简介' }],
          initialValue: currentRow.desc ? currentRow.desc : '',
        })(<TextArea placeholder="世界最好的域名" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="购入成本">
        {form.getFieldDecorator('cost', {
          initialValue: currentRow.cost ? currentRow.cost : 0,
        })(<InputNumber min={0} placeholder="1" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="续费成本">
        {form.getFieldDecorator('renew', {
          initialValue: currentRow.renew ? currentRow.renew : 0,
        })(<InputNumber min={0} placeholder="1" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="购入时间">
        {form.getFieldDecorator('buy', {
          initialValue: currentRow.buy ? moment(currentRow.buy) : moment(Date.now()),
        })(<DatePicker />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="注册商">
        {form.getFieldDecorator('registrar', {
          initialValue: currentRow.registrar ? currentRow.registrar : '',
        })(<Input placeholder="Alibaba Cloud" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="注册时间">
        {form.getFieldDecorator('create', {
          initialValue: currentRow.create ? moment(currentRow.create) : moment(Date.now()),
        })(<DatePicker />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="到期时间">
        {form.getFieldDecorator('expire', {
          initialValue: currentRow.expire ? moment(currentRow.expire) : moment(Date.now()),
        })(<DatePicker />)}
      </FormItem>
    </Modal>
  );
};

export default Form.create<CreateFormProps>()(CreateForm);
