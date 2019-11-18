import { Form, Modal } from 'antd';

import { FormComponentProps } from 'antd/es/form';
import React from 'react';
import TextArea from 'antd/lib/input/TextArea';
import moment from 'moment';

const FormItem = Form.Item;

interface ImportFormProps extends FormComponentProps {
  panel: any;
  importModalVisible: boolean;
  handleImport: (fieldsValue: any) => void;
  handleImportModalVisible: () => void;
}

const ImportForm: React.FC<ImportFormProps> = props => {
  const { importModalVisible, form, panel, handleImport, handleImportModalVisible } = props;

  const parseDomains = (text: string) => {
    let cat: any = null;
    const cats = [];
    text.split('\n').forEach((line, i) => {
      if (i === 0 && !line.startsWith('#')) {
        throw new Error(`第「${i}」行：必须以分类开头`);
      }
      if (line.startsWith('#')) {
        const catObj = line.substring(1).split(',');
        if (catObj.length !== 2) {
          throw new Error(`第「${i}」行：分类格式错误（#中文,英文）`);
        }
        if (cat) {
          cats.push(Object.assign({ name: '', name_en: '', domains: [] }, cat));
        }
        cat = { name: catObj[0], name_en: catObj[1], domains: [] };
      } else {
        const domainObj = line.split(',');
        if (domainObj.length !== 5) {
          throw new Error(
            `第「${i}」行：域名格式错误「域名,购入成本（可留空）,购入时间（可留空）,续费成本（可留空）,简介（可留空）」`,
          );
        }
        cat.domains.push({
          domain: domainObj[0],
          cost: parseInt(domainObj[1], 10),
          buy: moment(domainObj[2]),
          renew: parseInt(domainObj[3], 10),
          desc: domainObj[4],
        });
      }
    });
    cats.push(cat);
    return cats;
  };

  const okHandle = () => {
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      console.log(parseDomains(fieldsValue.text));
      handleImport(fieldsValue);
    });
  };

  return (
    <Modal
      destroyOnClose
      title="批量导入"
      visible={importModalVisible}
      onOk={okHandle}
      onCancel={() => handleImportModalVisible()}
    >
      <p style={{ fontSize: '14px' }}>
        ============== 导入格式 ==============
        <br />
        #分类中文,分类英文
        <br /> 域名,购入成本（可留空）,购入时间（可留空）,续费成本（可留空）,简介（可留空）
        <b>「留空可以，逗号不能省」</b>
        <br /> #单字符,Single Char
        <br /> qq.com,1000000,2005-01-01,69,腾讯网
        <br /> #双拼,Double Pinyin
        <br /> taobao.com,1000000,2006-01-01,69,淘宝网
      </p>
      {form.getFieldDecorator('panel_id', {
        initialValue: panel.id ? panel.id : 0,
      })(<input type="hidden" />)}
      <FormItem>
        {form.getFieldDecorator('text', {
          rules: [
            { required: true, message: '请输入导入的域名' },
            {
              validator: async value => {
                parseDomains(value);
              },
            },
          ],
        })(<TextArea />)}
      </FormItem>
    </Modal>
  );
};

export default Form.create<ImportFormProps>()(ImportForm);
