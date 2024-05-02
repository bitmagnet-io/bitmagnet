import { HumanTimePipe } from './human-time.pipe';

describe('HumanTimePipe', () => {
  it('create an instance', () => {
    const pipe = new HumanTimePipe();
    expect(pipe).toBeTruthy();
  });
});
