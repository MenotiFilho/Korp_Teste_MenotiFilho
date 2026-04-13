import { Component, Input, Output, EventEmitter, ContentChild, TemplateRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatSelectModule } from '@angular/material/select';

export interface FilterOption {
  value: string;
  label: string;
}

@Component({
  selector: 'app-data-table',
  imports: [FormsModule, MatIconModule, MatPaginatorModule, MatSelectModule],
  templateUrl: './data-table.component.html',
  styleUrl: './data-table.component.scss',
})
export class DataTableComponent {
  @Input() searchPlaceholder = '';
  @Input() searchWidth = 320;
  @Input() filterOptions: FilterOption[] = [];
  @Input() filterPlaceholder = 'Todos';
  @Input() totalItems = 0;
  @Input() pageSize = 10;

  @Output() searchChange = new EventEmitter<string>();
  @Output() filterChange = new EventEmitter<string>();
  @Output() pageChange = new EventEmitter<PageEvent>();

  searchTerm = '';
  filterValue = '';

  onSearch(): void {
    this.searchChange.emit(this.searchTerm);
  }

  onFilter(value: string): void {
    this.filterValue = value;
    this.filterChange.emit(value);
  }

  onPage(event: PageEvent): void {
    this.pageChange.emit(event);
  }
}
