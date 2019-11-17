import { Form, Input, Modal, InputNumber } from 'antd';

import { FormComponentProps } from 'antd/es/form';
import React from 'react';

const FormItem = Form.Item;

interface CreateFormProps extends FormComponentProps {
  isEdit: boolean;
  modalVisible: boolean;
  currentRow: any;
  handleAdd: (fieldsValue: any, isEdit: boolean) => void;
  handleModalVisible: () => void;
}

const CreateForm: React.FC<CreateFormProps> = props => {
  const { modalVisible, isEdit, currentRow, form, handleAdd, handleModalVisible } = props;

  const okHandle = () => {
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      Object.keys(fieldsValue).forEach(k => {
        if (typeof fieldsValue[k] === 'object') {
          fieldsValue[k] = fieldsValue[k].file;
        }
      });
      handleAdd(fieldsValue, isEdit);
    });
  };

  return (
    <Modal
      destroyOnClose
      title={`${isEdit ? '修改' : '新建'}分类`}
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
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="分类名">
        {form.getFieldDecorator('name', {
          rules: [{ required: true, message: '请输入分类名' }],
          initialValue: currentRow.name ? currentRow.name : '',
        })(<Input placeholder="国别" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="分类名「英」">
        {form.getFieldDecorator('name_en', {
          rules: [{ required: true, message: '请输入分类名「英」' }],
          initialValue: currentRow.name_en ? currentRow.name_en : '',
        })(<Input placeholder="ccTLD" />)}
      </FormItem>
      <FormItem labelCol={{ span: 6 }} wrapperCol={{ span: 15 }} label="排序">
        {form.getFieldDecorator('index', {
          initialValue: currentRow.index ? currentRow.index : 0,
        })(<InputNumber type="number" />)}
      </FormItem>
    </Modal>
  );
};

export default Form.create<CreateFormProps>()(CreateForm);
