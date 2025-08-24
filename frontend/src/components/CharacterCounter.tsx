import React from 'react';
import { CharacterCounterProps } from '../types/auth';

const CharacterCounter: React.FC<CharacterCounterProps> = ({
  current,
  max,
  showWarning = false,
  warningThreshold = 0.9
}) => {
  const warningClass = showWarning && current > max * warningThreshold ? 'warning' : '';

  return (
    <div className={`user-edit-character-count-internal ${warningClass}`}>
      <span className="current-count">{current}</span>/{max}
    </div>
  );
};

export default CharacterCounter;