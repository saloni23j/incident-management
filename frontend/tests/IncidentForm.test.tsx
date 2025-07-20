import { render, screen, fireEvent } from "@testing-library/react";
import IncidentForm from "../src/components/IncidentForm";
import { describe, expect, it, vi } from "vitest";

vi.mock("../src/api", () => ({
  createIncident: vi.fn(),
}));

describe("IncidentForm", () => {
  it("submits title and description", () => {
    const onIncidentCreated = vi.fn();
    render(<IncidentForm onIncidentCreated={onIncidentCreated} />);

    fireEvent.change(screen.getByPlaceholderText(/title/i), {
      target: { value: "Test Incident" },
    });
    fireEvent.change(screen.getByPlaceholderText(/description/i), {
      target: { value: "This is a test." },
    });
    fireEvent.click(screen.getByText(/submit/i));

    expect(onIncidentCreated).toHaveBeenCalledTimes(1);
  });
});
