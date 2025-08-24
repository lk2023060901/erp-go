import React from 'react';
import { Form, Input, Select, DatePicker } from 'antd';
import { TabProps } from '../types/auth';
import CharacterCounter from './CharacterCounter';

const { Option } = Select;
const { TextArea } = Input;

interface DetailInfoTabProps extends TabProps {
  // form instance not used in this component
}

const DetailInfoTab: React.FC<DetailInfoTabProps> = ({ 
  initialData, 
  onDataChange 
}) => {
  // 处理数据变化
  const handleDataChange = (field: string, value: any) => {
    const newData = { [field]: value };
    onDataChange?.(newData);
  };

  // 处理个人简介变化
  const handleBiographyChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    handleDataChange('biography', value);
  };

  return (
    <div>
      <div className="user-edit-form-row">
        <div className="user-edit-form-group">
          <label className="user-edit-form-label">性别</label>
          <Form.Item name="gender">
            <Select
              placeholder="请选择"
              allowClear
              onChange={(value) => handleDataChange('gender', value)}
            >
              <Option value="M">男</Option>
              <Option value="F">女</Option>
              <Option value="O">不愿透露</Option>
            </Select>
          </Form.Item>
        </div>

        <div className="user-edit-form-group">
          <label className="user-edit-form-label">生日</label>
          <Form.Item name="birth_date">
            <DatePicker
              className="user-edit-form-input"
              style={{ width: '100%' }}
              placeholder="请选择生日"
              format="YYYY-MM-DD"
              onChange={(_date, dateString) => handleDataChange('birth_date', dateString)}
            />
          </Form.Item>
        </div>
      </div>

      <div className="user-edit-form-row">
        <div className="user-edit-form-group">
          <label className="user-edit-form-label">手机号码</label>
          <Form.Item
            name="phone"
            rules={[
              { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码' }
            ]}
          >
            <Input
              className="user-edit-form-input"
              placeholder="请输入手机号码"
              onChange={(e) => handleDataChange('phone', e.target.value)}
            />
          </Form.Item>
        </div>

        <div className="user-edit-form-group">
          <label className="user-edit-form-label">联系地址</label>
          <Form.Item name="address">
            <Input
              className="user-edit-form-input"
              placeholder="详细地址"
              onChange={(e) => handleDataChange('address', e.target.value)}
            />
          </Form.Item>
        </div>
      </div>

      <div className="user-edit-form-row">
        <div className="user-edit-form-group">
          <label className="user-edit-form-label">紧急联系人</label>
          <Form.Item name="emergency_contact">
            <Input
              className="user-edit-form-input"
              placeholder="紧急联系人姓名"
              onChange={(e) => handleDataChange('emergency_contact', e.target.value)}
            />
          </Form.Item>
        </div>

        <div className="user-edit-form-group">
          <label className="user-edit-form-label">紧急联系电话</label>
          <Form.Item
            name="emergency_phone"
            rules={[
              { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码' }
            ]}
          >
            <Input
              className="user-edit-form-input"
              placeholder="紧急联系人电话"
              onChange={(e) => handleDataChange('emergency_phone', e.target.value)}
            />
          </Form.Item>
        </div>
      </div>

      <div className="user-edit-form-row">
        <div className="user-edit-form-group full-width">
          <label className="user-edit-form-label">个人简介</label>
          <Form.Item name="biography">
            <div className="user-edit-textarea-wrapper">
              <TextArea
                className="user-edit-form-textarea"
                placeholder="请输入个人简介..."
                maxLength={1000}
                rows={4}
                onChange={handleBiographyChange}
                defaultValue={initialData?.biography}
              />
              <CharacterCounter
                current={initialData?.biography?.length || 0}
                max={1000}
              />
            </div>
          </Form.Item>
        </div>
      </div>
    </div>
  );
};

export default DetailInfoTab;