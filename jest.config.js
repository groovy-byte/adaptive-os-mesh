module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  moduleNameMapper: {
    '^@voxel51/voodo$': '<rootDir>/__mocks__/voodoMock.tsx',
  },
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  testMatch: ['**/*.test.ts', '**/*.test.tsx'],
};
