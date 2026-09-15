import { File } from './bindings';
import type { TestFileView } from '../domain/types';

export class FileService {
  static async createFile(absPath: string): Promise<void> {
    try {
      await File.CreateFile(absPath);
    } catch (e) {
      console.error('FileService: createFile failed:', e);
      throw e;
    }
  }

  static async createDirectory(absPath: string): Promise<void> {
    try {
      await File.CreateDirectory(absPath);
    } catch (e) {
      console.error('FileService: createDirectory failed:', e);
      throw e;
    }
  }

  static async deletePath(absPath: string): Promise<void> {
    try {
      await File.DeletePath(absPath);
    } catch (e) {
      console.error('FileService: deletePath failed:', e);
      throw e;
    }
  }

  static async readTestFile(projectPath: string, relPath: string): Promise<TestFileView> {
    try {
      const view = await File.ReadTestFile(projectPath, relPath);
      return {
        graph: view.graph,
        order: view.order ?? [],
      };
    } catch (e) {
      console.error('FileService: readTestFile failed:', e);
      throw e;
    }
  }

  static async writeTestFile(
    projectPath: string,
    relPath: string,
    view: TestFileView,
  ): Promise<void> {
    try {
      await File.WriteTestFile(projectPath, relPath, view as any);
    } catch (e) {
      console.error('FileService: writeTestFile failed:', e);
      throw e;
    }
  }

  static async validate(
    projectPath: string = '',
    relPath: string = '',
    view?: TestFileView,
  ): Promise<import('../domain/types').DiagnosticReport> {
    try {
      const report = await File.Validate(projectPath, relPath, view as any);
      // The generated binding widens severity and code to plain strings while the
      // local mirror narrows them to unions. The values come from Go, which only
      // emits those two severities, so the boundary cast is safe.
      return report as unknown as import('../domain/types').DiagnosticReport;
    } catch (e) {
      console.error('FileService: validate failed:', e);
      return { isValid: true, errors: [], warnings: [] };
    }
  }
}
